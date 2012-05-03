// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package kindi

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"strings"
	"testing"

	"github.com/uwedeportivo/shared/util"
)

const (
	metadata = "some metadata"
	payload  = "The quick brown fox jumped over the lazy programmer."
)

func checkRecipient(t *testing.T, recipient *Identity, r io.Reader) {
	cs, err := recipient.DecryptCipherStream(r)
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	abtr := util.NewAllButTailReader(r, cs.Hash.Size(), 512)

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

func TestEncryption(t *testing.T) {
	size := 1024
	identities := make([]*rsa.PrivateKey, 3)
	for i := 0; i < len(identities); i++ {
		identity, err := rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			t.Fatalf("failed %v", err)
		}
		identities[i] = identity
	}

	publicKeys := make([]*rsa.PublicKey, len(identities))
	for i := 0; i < len(identities); i++ {
		publicKeys[i] = &identities[i].PublicKey
	}

	r := strings.NewReader(payload)
	bufOne := new(bytes.Buffer)
	bufTwo := new(bytes.Buffer)

	sender := &Identity{
		Email:      "sender@testing.com",
		PrivateKey: identities[0],
	}

	recipientOne := &Identity{
		Email:      "recipientOne@testing.com",
		PrivateKey: identities[1],
	}

	recipientTwo := &Identity{
		Email:      "recipientTwo@testing.com",
		PrivateKey: identities[2],
	}

	err := sender.Encrypt(io.MultiWriter(bufOne, bufTwo), []byte(metadata), r, publicKeys[1:])
	if err != nil {
		t.Fatalf("failed %v", err)
	}

	checkRecipient(t, recipientOne, bufOne)
	checkRecipient(t, recipientTwo, bufTwo)
}
