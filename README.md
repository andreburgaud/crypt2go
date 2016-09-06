# Crypt2go

Some basic Golang packages complementing existing standard library `crypto` packages and the extension packages `x/crypto`.

## Installation

```
$ go get github/andreburgaud/crypt2go
```

or to update to the latest version:

```
$ go get -u github/andreburgaud/crypt2go
```

## Disclaimer

I'm, by no mean, an expert in cryptography and welcome any comment or suggestion to improve the code included in this repository.

## ECB (Electronic Codebook)

The ECB mode of operation should **NOT** be used anymore. This code was written to facilate migrating legacy data encrypted with ECB.

There is plenty of literature explaining why ECB should not be used, starting with https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation, section *Electronic Codebook (ECB)*.

Nevertheless, if like me, someone needed to solve a problem with legacy software and use ECB, this code might be helpful.

### Blowfish encryption in ECB mode with padding

```go
package main

import (
  "fmt"

  "golang.org/x/crypto/blowfish"

  "github.com/andreburgaud/crypt2go/ecb"
  "github.com/andreburgaud/crypt2go/padding"
)

func encrypt(pt, key []byte) []byte {
  block, err := blowfish.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }
  mode := ecb.NewECBEncrypter(block)
  padder := padding.NewPkcs5Padding()
  pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
  if err != nil {
    panic(err.Error())
  }
  ct := make([]byte, len(pt))
  mode.CryptBlocks(ct, pt)
  return ct
}

func decrypt(ct, key []byte) []byte {
  block, err := blowfish.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }
  mode := ecb.NewECBDecrypter(block)
  pt := make([]byte, len(ct))
  mode.CryptBlocks(pt, ct)
  padder := padding.NewPkcs5Padding()
  pt, err = padder.Unpad(pt) // unpad plaintext after decryption
  if err != nil {
    panic(err.Error())
  }
  return pt
}

func example() {
  pt := []byte("Some plain text")
  key := []byte("a_very_secret_key")

  ct := encrypt(pt, key)
  fmt.Printf("Ciphertext: %x\n", ct)

  recovered_pt := decrypt(ct, key)
  fmt.Printf("Recovered plaintext: %s\n", recovered_pt)
}

func main() {
  example()
}
```

### AES Encryption in ECB mode and with padding

```go
import (
  "fmt"

  "crypto/aes"

  "github.com/andreburgaud/crypt2go/ecb"
  "github.com/andreburgaud/crypt2go/padding"
)

func encrypt(pt, key []byte) []byte {
  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }
  mode := ecb.NewECBEncrypter(block)
  padder := padding.NewPkcs7Padding(mode.BlockSize())
  pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
  if err != nil {
    panic(err.Error())
  }
  ct := make([]byte, len(pt))
  mode.CryptBlocks(ct, pt)
  return ct
}

func decrypt(ct, key []byte) []byte {
  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }
  mode := ecb.NewECBDecrypter(block)
  pt := make([]byte, len(ct))
  mode.CryptBlocks(pt, ct)
  padder := padding.NewPkcs7Padding(mode.BlockSize())
  pt, err = padder.Unpad(pt) // unpad plaintext after decryption
  if err != nil {
    panic(err.Error())
  }
  return pt
}

func example() {
  pt := []byte("Some plain text")
  // Key size for AES is either: 16 bytes (128 bits), 24 bytes (192 bits) or 32 bytes (256 bits)
  key := []byte("secretkey16bytes")

  ct := encrypt(pt, key)
  fmt.Printf("Ciphertext: %x\n", ct)

  recovered_pt := decrypt(ct, key)
  fmt.Printf("Recovered plaintext: %s\n", recovered_pt)
}

func main() {
  fmt.Println("AES encryption with ECB and PKCS7 padding")
  example()
}
```

## Padding

Both ECB (Electronic Codebook) and CBC (Cipher Block Chaining) require blocks of fixed size. In order to comply with this requirement, it is necessary to `pad` the plain text to a size multiple of the block size in order to perform any encryption with these modes of operation.

The `padding` package exposes simple functions to provide a way to `pad` and `unpad` a given plaintext respectively prior to encryption and after decryption.

The code examples in the previous section shows encryption examples with AES, Blowfish in ECB mode. Blowfish encrypts blocks of 8 bytes hence using the padding type described in the https://tools.ietf.org/html/rfc2898 *PKCS #5: Password-Based Cryptography Specification Version 2.0*. Whereas AES requires blocks of 16 bytes (128 bits). The padding type in the second example is based on https://tools.ietf.org/html/rfc2315 *PKCS #7: Cryptographic Message Syntax Version 1.5*.

The only difference between the two specs is only that PKCS #5 accommodates only for blocks of 8 bytes. The `padding` package reflects that and exposes two builders, resepectively `NewPkcs5Padding()` that embeds an hard-coded value for a block size of 8, while `NewPkcs7Padding(int blockSize)` takes a parameter for the block size. Nothing prevents to use `NewPkcs7Padding` with a block size of 8 to work with an encryption scheme working on blocks of 8 bytes, like *Blowfish*.

## Additional Details

Check the unit tests in the code and the godoc string for additional details.

## License

Some of the code, in particular the `ecb` package is directly modeled after the Golang code (`cipher/cbc.go`). To avoid any license conflicts, `crypt2go` is released under a BSD license. See the LICENSE file in the repository.

The Go Authors license is available at https://golang.org/LICENSE. The reference to this license is made in the code where appropriate.


