package signer

import (
	"testing"
)

func TestSigner(t *testing.T) {
	signer := NewSigner("secret")

	t.Run("Test signer with valid value", func(t *testing.T) {
		signed := signer.Sign("20")
		valid := signer.Validate("20", signed)
		expected := true

		if valid != expected {
			t.Errorf("Expected %v, got %v", expected, valid)
		}
	})

	t.Run("Test signer with invalid value", func(t *testing.T) {
		valid := signer.Validate("20", "fakedata")
		expected := false

		if valid != expected {
			t.Errorf("Expected %v, got %v", expected, valid)
		}
	})
}
