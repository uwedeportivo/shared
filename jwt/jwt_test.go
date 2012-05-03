// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package jwt

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	issued := time.Unix(1333488524, 0)
	expected, err := ioutil.ReadFile("testdata/example.token")
	if err != nil {
		t.Fatalf("error reading expected value: %v", err)
	}

	request := map[string]string{
		"name":         "Piece of Cake",
		"description":  "Virtual chocolate cake to fill your virtual tummy",
		"price":        "10.50",
		"currencyCode": "USD",
		"sellerData":   "user_id:1224245,offer_code:3098576987,affiliate:aksdfbovu9j",
	}

	token := Token{
		Request:          request,
		Response:         nil,
		Issued:           issued,
		Expires:          issued.Add(time.Hour),
		SellerIdentifier: "1337133713371337",
		SellerSecret:     "PWGknVgi6zt_BU1qrO1h",
	}

	jot, err := Encode(token)
	if err != nil {
		t.Fatalf("error generating JWT token: %v", err)
	}

	if jot != string(expected) {
		t.Fatalf("generated JWT token differs from expected value")
	}
}

func TestDecode(t *testing.T) {
	jot, err := ioutil.ReadFile("testdata/example.token")
	if err != nil {
		t.Fatalf("error reading jot value: %v", err)
	}

	token, err := Decode(string(jot), "1337133713371337", "PWGknVgi6zt_BU1qrO1h", true)
	if err != nil {
		t.Fatalf("error decoding jot: %v", err)
	}

	issued := time.Unix(1333488524, 0)
	if issued != token.Issued {
		t.Fatalf("unexpected jot issued time: %v", token.Issued)
	}

	expires := issued.Add(time.Hour)
	if expires != token.Expires {
		t.Fatalf("unexpected jot expires time: %v", token.Expires)
	}

	if "Piece of Cake" != token.Request["name"] {
		t.Fatalf("unexpected jot request name: %v", token.Request["name"])
	}
	if "Virtual chocolate cake to fill your virtual tummy" != token.Request["description"] {
		t.Fatalf("unexpected jot request description: %v", token.Request["description"])
	}
	if "10.50" != token.Request["price"] {
		t.Fatalf("unexpected jot request price: %v", token.Request["price"])
	}
	if "USD" != token.Request["currencyCode"] {
		t.Fatalf("unexpected jot request currencyCode: %v", token.Request["currencyCode"])
	}
	if "user_id:1224245,offer_code:3098576987,affiliate:aksdfbovu9j" != token.Request["sellerData"] {
		t.Fatalf("unexpected jot request sellerData: %v", token.Request["sellerData"])
	}
}

func TestDecodePadding(t *testing.T) {
	jot, err := ioutil.ReadFile("testdata/google.token")
	if err != nil {
		fmt.Printf("error reading jot value: %v\n", err)
	}

	_, err = Decode(string(jot), "16152247716018635319", "YhdfG5Pr6c8X3EIXd-jg9w", false)
	if err != nil {
		t.Fatalf("error decoding jot: %v", err)
	}
}

