package ecb

import (
	"crypto/aes"
	"encoding/hex"
	"testing"
)

func TestECBEncrypterAESPanic(t *testing.T) {
	// if the encryption does not panic the test will fail
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("input plaintext was too short expected to panic")
		}
	}()
	key := []byte("AES256Key-32Characters1234567890")
	plaintext := []byte("exampleplaintex") // input is too short not a multiple of 16, should be padded, therefore will panic
	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("NewCipher(%d bytes) = %s", len(key), err)
	}
	encrypter := NewECBEncrypter(c)
	ciphertext := make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)
}

func TestECBDecrypterAESPanic(t *testing.T) {
	// if the decryption does not panic the test will fail
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("input plaintext was too short expected to panic")
		}
	}()
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, err := hex.DecodeString("717FADE7B97198A8C2F67766FBAC7B") // correct value was "717FADE7B97198A8C2F67766FBAC7B07" (1 byte missing => should panic)

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("NewCipher(%d bytes) = %s", len(key), err)
	}
	encrypter := NewECBEncrypter(c)
	plaintext := make([]byte, len(ciphertext))
	encrypter.CryptBlocks(plaintext, ciphertext)
}
