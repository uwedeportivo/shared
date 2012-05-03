// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package kindi

import (
	"bytes"
	"codemanic.com/util"
	"crypto/rand"
	"io"
	"strings"

	"testing"
)

func TestSymmetricEncryption(t *testing.T) {
	metadata := "some metadata"
	payload := "The quick brown fox jumped over the lazy programmer."

	r := strings.NewReader(payload)
	buf := new(bytes.Buffer)

	symmetricKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, symmetricKey)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	iv := make([]byte, BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	// ---- encrypt -----
	cs, err := NewCipherStream(symmetricKey, iv)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	err = cs.EncryptMetadata(buf, []byte(metadata))
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	err = cs.EncryptPayload(buf, r)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	err = cs.AppendHash(buf)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	// ----- decrypt ------
	cs, err = NewCipherStream(symmetricKey, iv)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	abtr := util.NewAllButTailReader(buf, cs.Hash.Size(), 512)

	decryptedMetadata, err := cs.DecryptMetadata(abtr)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	decrypted := new(bytes.Buffer)

	err = cs.DecryptPayload(decrypted, abtr)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	if !bytes.Equal(abtr.Tail(), cs.Hash.Sum(nil)) {
		t.Fatalf("hash check failed")
	}

	if metadata != string(decryptedMetadata) {
		t.Fatalf("decrypted metadata %s different from expected metadata %s", string(decryptedMetadata), metadata)
	}

	if payload != string(decrypted.Bytes()) {
		t.Fatalf("decrypted payload %s different from expected payload %s", string(decrypted.Bytes()), payload)
	}
}
