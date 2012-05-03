// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package kindi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"io"

	"github.com/uwedeportivo/shared/util"
)

type Identity struct {
	Email      string
	PrivateKey *rsa.PrivateKey
}

type JSONKindiCertificate struct {
	Email string
	Bytes []byte
}

func (me *Identity) Sign() ([]byte, error) {
	hash := sha1.New()
	hash.Write([]byte(me.Email))
	sum := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, me.PrivateKey, crypto.SHA1, sum)
}

func VerifySignature(senderEmail string, senderSignature []byte, senderCerts []*rsa.PublicKey) error {
	hash := sha1.New()
	hash.Write([]byte(senderEmail))
	sum := hash.Sum(nil)

	for _, cert := range senderCerts {
		err := rsa.VerifyPKCS1v15(cert, crypto.SHA1, sum, senderSignature)
		if err == nil {
			return nil
		}
	}
	return errors.New("sender signature did not verify")
}

func (me *Identity) DecryptCipherStream(r io.Reader) (*CipherStream, error) {
	var n int64
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}

	var decrypted []byte
	var found bool

	for i := 0; i < int(n); i++ {
		encryptedSymmetricKey, err := util.ReadLengthEncoded(r)
		if err != nil {
			return nil, err
		}
		if !found {
			hash := sha1.New()
			decrypted, err = rsa.DecryptOAEP(hash, rand.Reader, me.PrivateKey, encryptedSymmetricKey, nil)
			if err == nil {
				found = true
			}
		}
	}

	if !found {
		return nil, errors.New("could not decipher symmetric key")
	}

	symmetricKey := decrypted[0:32]
	iv := decrypted[32:]

	return NewCipherStream(symmetricKey, iv)
}

func (me *Identity) Encrypt(w io.Writer, metadata []byte, payload io.Reader, recipients []*rsa.PublicKey) error {
	if len(recipients) == 0 {
		return errors.New("no recipients")
	}

	symmetricKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, symmetricKey)
	if err != nil {
		return err
	}

	iv := make([]byte, BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	err = encryptSymmetricKey(w, symmetricKey, iv, recipients)
	if err != nil {
		return err
	}

	cs, err := NewCipherStream(symmetricKey, iv)
	if err != nil {
		return err
	}

	err = cs.EncryptMetadata(w, metadata)
	if err != nil {
		return err
	}

	err = cs.EncryptPayload(w, payload)
	if err != nil {
		return err
	}

	return cs.AppendHash(w)
}

func encryptSymmetricKey(w io.Writer, symmetricKey []byte, iv []byte, recipients []*rsa.PublicKey) error {
	if len(recipients) == 0 {
		return errors.New("no recipients")
	}

	err := binary.Write(w, binary.BigEndian, int64(len(recipients)))
	if err != nil {
		return err
	}

	payload := make([]byte, len(symmetricKey)+len(iv))
	copy(payload[0:len(symmetricKey)], symmetricKey)
	copy(payload[len(symmetricKey):], iv)

	for _, recipient := range recipients {
		b, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, recipient, payload, nil)
		if err != nil {
			return err
		}
		err = util.WriteLengthEncoded(w, b)
		if err != nil {
			return err
		}
	}
	return nil
}
