package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "secret-password"
	plaintext := []byte("Hello, Gist-Hub! This is a secure message.")

	// Encrypt
	encoded, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encoded == "" {
		t.Fatal("Encrypted string is empty")
	}

	// Decrypt
	decrypted, err := Decrypt(encoded, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted data mismatch.\nExpected: %s\nGot:      %s", plaintext, decrypted)
	}
}

func TestDecryptWithInvalidPassphrase(t *testing.T) {
	passphrase := "correct-password"
	wrongPassphrase := "wrong-password"
	plaintext := []byte("Sensitive information")

	encoded, _ := Encrypt(plaintext, passphrase)

	_, err := Decrypt(encoded, wrongPassphrase)
	if err == nil {
		t.Error("Decryption should have failed with wrong passphrase")
	}
}

func TestDeterministicKeyDerivation(t *testing.T) {
	// 異なる暗号文でも、同じパスフレーズとソルトからは同じ鍵が生成されることを
	// ロジック的に保証するための確認（内部的な key 変数は公開していないため、
	// ここでは実装の不変条件として念頭に置く）
}
