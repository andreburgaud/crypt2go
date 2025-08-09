# Crypt2go

[![GoDoc](https://godoc.org/github.com/andreburgaud/crypt2go?status.svg)](https://godoc.org/github.com/andreburgaud/crypt2go)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreburgaud/crypt2go)](https://goreportcard.com/report/github.com/andreburgaud/crypt2go)

**Crypt2go** includes [Go](https://go.dev/) packages complementing existing standard library [crypto](https://pkg.go.dev/crypto) packages and the extension packages [x/crypto](https://pkg.go.dev/golang.org/x/crypto).

## Installation

```
go get github.com/andreburgaud/crypt2go
```

Or to update to the latest version:

```
go get -u github.com/andreburgaud/crypt2go
```

## Development

In the `crypt2go` directory, execute `make` or `make help` to display the build commands.

### Test

```
make test
```

### Run Examples

To execute the code examples found in this document (`README`):

```
make run
```

## Disclaimer

I'm not an expert in cryptography, and I welcome any comments or suggestions to improve the code included in this repository.

## ECB (Electronic Codebook)

The [ECB mode of operation](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_codebook_(ECB)) should **NOT** be used anymore. I wrote **Crypt2go** to facilitate migrating legacy data encrypted with ECB. Plenty of literature explains why ECB is not recommended in cryptographic protocols. Nevertheless, this code might be helpful if someone needs to solve a problem with legacy software using ECB.

## Examples

### AES Encryption in ECB mode with padding

```go
package main

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
  pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
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
  // Key size for AES is either: 16 bytes (128 bits), 24 bytes (192 bits), or 32 bytes (256 bits)
  key := []byte("secretkey16bytes")

  ct := encrypt(pt, key)
  fmt.Printf("Ciphertext: %x\n", ct)

  recoveredPt := decrypt(ct, key)
  fmt.Printf("Recovered plaintext: %s\n", recoveredPt)
}

func main() {
  fmt.Println("AES encryption with ECB and PKCS7 padding")
  example()
}
```

### Blowfish encryption in ECB mode with padding

**Note**: The Golang blowfish package is deprecated and should not be used for any new system (https://pkg.go.dev/golang.org/x/crypto/blowfish). Nevertheless, you could still benefit from the following example if you needed to convert legacy encryption from blowfish with ECB to a secure encryption algorithm and mode of encryption, for example, AES with GCM ([Gallois Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode)). To allow linting via golangci-lint to complete successfully, the import in the blowfish example is intentionally configured not to lint.

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
  pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
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

  recoveredPt := decrypt(ct, key)
  fmt.Printf("Recovered plaintext: %s\n", recoveredPt)
}

func main() {
  fmt.Println("Blowfish encryption with ECB and PKCS5 padding")
  example()
}
```


## Padding

ECB (Electronic Codebook) and CBC (Cipher Block Chaining) require fixed-size blocks. Encryption with ECB and CBC modes requires padding the plain text to a size multiple of the block size.

The `padding` package exposes simple functions to provide a way to `pad` before encryption and `unpad` a given plaintext after decryption.

The code examples in the previous sections show encryption patterns with [Blowfish](https://en.wikipedia.org/wiki/Blowfish_(cipher)) and [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) in ECB mode. Blowfish encrypts blocks of 8 bytes (64 bits) using the padding type described in the [PKCS #5: Password-Based Cryptography Specification Version 2.1](https://tools.ietf.org/html/rfc8018). In contrast, AES requires blocks of 16 bytes (128 bits). The padding type in the second example follows the specs documented in [PKCS #7: Cryptographic Message Syntax Version 1.5](https://tools.ietf.org/html/rfc2315).

The only difference between the two specs is that PKCS #5 accommodates only blocks of 8 bytes. The `padding` package reflects that and exposes two builders, respectively `NewPkcs5Padding()` that embeds a hard-coded value for a block size of 8, while `NewPkcs7Padding(int blockSize)` takes a parameter for the block size. Nothing prevents using `NewPkcs7Padding` with a block size of 8 to work with an encryption scheme on blocks of 8 bytes, like *Blowfish*.

### Full block of padding

Padding is always performed to ensure there is no ambiguity for the receiver of a message, even if the message is of an exact multiple-block size. It is intentional and complies with the [NIST Recommendation for Block Cipher Modes of Operation](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf) **Appendix A: Padding** (page 17), and the following RFCs:

* [PKCS #5: Password-Based Cryptography Specification Version 2.1](https://tools.ietf.org/html/rfc8018)
* [Privacy Enhancement for Internet Electronic Mail: Part III: Algorithms, Modes, and Identifiers](https://tools.ietf.org/html/rfc1423)

The padding goes as follows for blocks of 8 bytes (8 octets):

Given a message `M`, we obtain an encoded message `EM` by concatenating `M` with a padding string `PS`:

```
EM = M || PS
```

The padding string `PS` consists of `8 - (||M|| mod 8)` octets, each with value `8 - (||M|| mod 8)`. Examples:

```
PS = 01, if ||M|| mod 8 = 7
PS = 02 02, if ||M|| mod 8 = 6
...
PS = 08 08 08 08 08 08 08 08, if ||M|| mod 8 = 0
```

The last example is essential. Yes, it will intentionally add an entire padding block. Doing so removes the ambiguity for the receiver, that expects every message to be padded.

To illustrate what would happen if some messages are not padded, let's take an example of a message with the last octet with the value `01`. As the receiver of this message, should I remove the padding `01`? Or, is the last byte `01` part of a message not padded because it was an exact multiple-block size?

If the receiver knows that every message is padded, even if this results in a message padded with a whole block of `08`, there is no ambiguity.

Another approach to removing this ambiguity would be to provide a separate indicator that would remove this ambiguity. An example would be to give a message length indicator.

The implementation in this package relies on padding every message (see method `padding.Pad()`).

## Additional Examples

See the unit tests or the example tests in the respective package directories.

## License

The **Crypt2go** `ecb` package is directly modeled after the [CBC Go code](https://go.dev/src/crypto/cipher/cbc.go) released under a [BSD license](https://go.dev/LICENSE). To avoid license conflicts, **Crypt2go** is also released under a BSD license.

See the [LICENSE](LICENSE.md) file in the repository.
