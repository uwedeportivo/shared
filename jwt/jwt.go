// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/uwedeportivo/shared/util"
)

type Token struct {
	Request          map[string]string
	Response         map[string]string
	Issued           time.Time
	Expires          time.Time
	SellerIdentifier string
	SellerSecret     string
}

func signature(headerAndBody string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, headerAndBody)
	sig := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sig)
}

func base64Decode(s string) ([]byte, error) {
	encoded := s
	if m := len(encoded) % 4; m != 0 {
  		encoded += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(encoded)
}

func verifyHeader(header string) error {
	headerJson, err := base64Decode(header)
	if err != nil {
		return err
	}

	var headerBlob interface{}
	err = json.Unmarshal(headerJson, &headerBlob)
	if err != nil {
		return err
	}

	typBlob := util.JsonPath(headerBlob, "typ")

	if typBlob != nil {
		typ := typBlob.(string)

		if typ != "JWT" {
			return errors.New("unexpected typ value in header: " + typ)
		}
	}	

	algBlob := util.JsonPath(headerBlob, "alg")
	alg := algBlob.(string)

	if alg != "HS256" {
		return errors.New("unexpected alg value in header: " + alg)
	}
	return nil
}

func processBody(body string, sellerIdentifier string, sellerSecret string) (*Token, error) {
	bodyJson, err := base64Decode(body)
	if err != nil {
		return nil, err
	}

	var bodyBlob interface{}
	err = json.Unmarshal(bodyJson, &bodyBlob)
	if err != nil {
		return nil, err
	}

	audBlob := util.JsonPath(bodyBlob, "aud")
	aud := audBlob.(string)

	if aud != "Google" && aud != sellerIdentifier {
		return nil, errors.New("unexpected aud value in body: " + aud)
	}

	typBlob := util.JsonPath(bodyBlob, "typ")
	typ := typBlob.(string)

	if typ != "google/payments/inapp/item/v1" && typ != "google/payments/inapp/item/v1/postback/buy" {
		return nil, errors.New("unexpected typ value in body: " + typ)
	}

	issBlob := util.JsonPath(bodyBlob, "iss")
	iss := issBlob.(string)

	if iss != sellerIdentifier && iss != "Google" {
		return nil, errors.New("unexpected iss value in body: " + iss)
	}

	iatBlob := util.JsonPath(bodyBlob, "iat")

	var iat int64
	switch i := iatBlob.(type) {
		case nil:
			err = errors.New("iat is missing")
		case int:
			iat = int64(i)
		case int64:
			iat = i
		case float64:
			iat = int64(i)		
		case string:
			iat, err = strconv.ParseInt(i, 10, 64)
		default:
			err = errors.New("iat has unexpected type")
	}
	if err != nil {
		return nil, err
	}

	expBlob := util.JsonPath(bodyBlob, "exp")

	var exp int64
	switch i := expBlob.(type) {
		case nil:
			err = errors.New("exp is missing")
		case int:
			exp = int64(i)
		case int64:
			exp = i	
		case float64:
			exp = int64(i)	
		case string:
			exp, err = strconv.ParseInt(i, 10, 64)
		default:
			err = errors.New("exp has unexpected type")
	}
	if err != nil {
		return nil, err
	}

	rv := new(Token)

	rv.Issued = time.Unix(iat, 0)
	rv.Expires = time.Unix(exp, 0)
	rv.SellerIdentifier = sellerIdentifier
	rv.SellerSecret = sellerSecret

	requestBlob := util.JsonPath(bodyBlob, "request")
	request := requestBlob.(map[string]interface{})

	rv.Request = make(map[string]string)

	for k, v := range request {
		vStr := v.(string)
		rv.Request[k] = vStr
	}

	responseBlob := util.JsonPath(bodyBlob, "response")
	if responseBlob != nil {
		response := responseBlob.(map[string]interface{})

		rv.Response = make(map[string]string)

		for k, v := range response {
			vStr := v.(string)
			rv.Response[k] = vStr
		}
	}
	return rv, nil
}

func Encode(token Token) (string, error) {
	header := base64.URLEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"HS256"}`))
	body := make(map[string]interface{})

	body["aud"] = "Google"
	body["typ"] = "google/payments/inapp/item/v1"
	body["iss"] = token.SellerIdentifier

	body["iat"] = strconv.FormatInt(token.Issued.Unix(), 10)
	body["exp"] = strconv.FormatInt(token.Expires.Unix(), 10)

	body["request"] = token.Request

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	payload := base64.URLEncoding.EncodeToString(bodyJson)

	headerAndBody := header + "." + payload
	sig := signature(headerAndBody, token.SellerSecret)

	return headerAndBody + "." + sig, nil
}

func Decode(jot string, sellerIdentifier string, sellerSecret string, verify bool) (*Token, error) {
	parts := strings.Split(jot, ".")

	if len(parts) != 3 {
		return nil, errors.New("jot invalid, missing parts")
	}

	err := verifyHeader(parts[0])
	if err != nil {
		return nil, err
	}

	if verify {
		sig := signature(parts[0]+"."+parts[1], sellerSecret)

		if sig != parts[2] {
			return nil, errors.New("hmac verification failed")
		}
	}	

	return processBody(parts[1], sellerIdentifier, sellerSecret)
}
