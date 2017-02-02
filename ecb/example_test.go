package ecb

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blowfish"
)

func ExampleNewECBEncrypter() {
	key := []byte("SomeBlowfishKey")
	plaintext := []byte("exampleplaintext")
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := NewECBEncrypter(block)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	fmt.Printf("%X\n", ciphertext)
	// Output: 66E6E737F23E570ACA5710BA3C0E321A
}

func ExampleNewECBEncrypter_second() {
	key := []byte("AES256Key-32Characters1234567890")
	plaintext := []byte("exampleplaintext")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := NewECBEncrypter(block)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	fmt.Printf("%X\n", ciphertext)
	// Output: 717FADE7B97198A8C2F67766FBAC7B07
}

func ExampleNewECBDecrypter() {
	key := []byte("SomeBlowfishKey")
	ciphertext, err := hex.DecodeString("66E6E737F23E570ACA5710BA3C0E321A")
	if err != nil {
		panic(err.Error())
	}
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := NewECBDecrypter(block)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	fmt.Printf("%s\n", string(plaintext))
	// Output: exampleplaintext
}

func ExampleNewECBDecrypter_second() {
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, err := hex.DecodeString("717FADE7B97198A8C2F67766FBAC7B07")
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := NewECBDecrypter(block)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	fmt.Printf("%s\n", string(plaintext))
	// Output: exampleplaintext
}
