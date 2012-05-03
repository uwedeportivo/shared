// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package kindi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	"github.com/uwedeportivo/shared/util"
)

const BlockSize = aes.BlockSize

type CipherStream struct {
	Stream cipher.Stream
	Hash   hash.Hash
}

func NewCipherStream(symmetricKey []byte, iv []byte) (*CipherStream, error) {
	c, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("Failed to create cipher")
	}

	stream := cipher.NewOFB(c, iv)

	if stream == nil {
		return nil, errors.New("Failed to create cipher.Stream")
	}

	hashSeed := make([]byte, 64)
	c.Encrypt(hashSeed, hashSeed)

	return &CipherStream{
		Stream: stream,
		Hash:   hmac.New(sha256.New, hashSeed),
	}, nil
}

func (cs *CipherStream) EncryptMetadata(w io.Writer, metadata []byte) error {
	encryptWriter := &cipher.StreamWriter{S: cs.Stream, W: io.MultiWriter(w, cs.Hash)}

	return util.WriteLengthEncoded(encryptWriter, metadata)
}

func (cs *CipherStream) EncryptPayload(w io.Writer, r io.Reader) error {
	encryptWriter := &cipher.StreamWriter{S: cs.Stream, W: io.MultiWriter(w, cs.Hash)}

	_, err := io.Copy(encryptWriter, r)
	return err
}

func (cs *CipherStream) AppendHash(w io.Writer) error {
	_, err := w.Write(cs.Hash.Sum(nil))
	return err
}

func (cs *CipherStream) DecryptMetadata(r io.Reader) ([]byte, error) {
	decryptReader := &cipher.StreamReader{S: cs.Stream, R: &util.HashReader{R: r, H: cs.Hash}}

	return util.ReadLengthEncoded(decryptReader)
}

func (cs *CipherStream) DecryptPayload(w io.Writer, r io.Reader) error {
	decryptReader := &cipher.StreamReader{S: cs.Stream, R: &util.HashReader{R: r, H: cs.Hash}}
	_, err := io.Copy(w, decryptReader)
	return err
}
