package encrypt

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

const prefix = "ENC*"

type Cipher interface {
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
}

var _ Cipher = (*aeadCipher)(nil)

type aeadCipher struct {
	aead cipher.AEAD
}

// New return new instance of ChaChaCiper
func New(key []byte) (*aeadCipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key length must me 32 byte")
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	return &aeadCipher{
		aead: aead,
	}, nil
}

// Decrypt decrypts input
func (c *aeadCipher) Decrypt(input []byte) ([]byte, error) {
	if !strings.HasPrefix(string(input), prefix) {
		return input, nil
	}

	// If we use base64.StdEncoding.Decode and base64.StdEncoding.DecodedLen
	// for decoding, it will return slightly bigger slice and fill the
	// rest with \x00 bytes. This invalidates the message and decryption failed
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(input[4:]))
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r)
	if err != nil {
		return []byte{}, err
	}

	if len(buf.Bytes()) < 24 {
		return []byte{}, fmt.Errorf("failed to decrypt message")
	}
	nonce := buf.Bytes()[0:24]
	data := buf.Bytes()[24:]

	plaintext, err := c.aead.Open(nil, nonce, data, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to decrypt or authenticate message: %s", err)
	}
	return plaintext, nil
}

// Encrypt encrypts input using 24 byte nonce
func (c *aeadCipher) Encrypt(input []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	rand.Read(nonce)
	encrypted := c.aead.Seal(nil, nonce, input, nil)

	cipherText := append(nonce, encrypted...)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(result, cipherText)

	// return prefix+base64(nonce+ciphertext)
	final := append([]byte(prefix), result...)
	return final, nil
}