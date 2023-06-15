package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Signer struct {
	secretKey string
}

func NewSigner(secretKey string) *Signer {
	return &Signer{
		secretKey: secretKey,
	}
}

func (s *Signer) Sign(value string) string {
	return signData(value, s.secretKey)
}

func (s *Signer) Validate(value, signedValue string) bool {
	return validateCookie(value, signedValue, s.secretKey)
}

func signData(value, secretKey string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(value))
	signedValue := hex.EncodeToString(hash.Sum(nil))

	return signedValue
}

func validateCookie(value, signedValue, secretKey string) bool {
	return hmac.Equal([]byte(signedValue), []byte(signData(value, secretKey)))
}
