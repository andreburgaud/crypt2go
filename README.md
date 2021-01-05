# Crypt2go

[![GoDoc](https://godoc.org/github.com/andreburgaud/crypt2go?status.svg)](https://godoc.org/github.com/andreburgaud/crypt2go)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreburgaud/crypt2go)](https://goreportcard.com/report/github.com/andreburgaud/crypt2go)

Some basic Golang packages complementing existing standard library `crypto` packages and the extension packages `x/crypto`.

## Installation

```
$ go get github.com/andreburgaud/crypt2go
```

or to update to the latest version:

```
$ go get -u github.com/andreburgaud/crypt2go
```

## Development

In the `crypt2go` directory, execute `make` or `make help` to display the build commands.

### Test

```
$ make test
```

### Run Examples

To execute the examples similar to these in the `README` (this file):

```
$ make run
```

## Disclaimer

I'm, by no means, an expert in cryptography and welcome any comments or suggestions to improve the code included in this repository.

## ECB (Electronic Codebook)

The ECB mode of operation should **NOT** be used anymore. This code was written to facilitate migrating legacy data encrypted with ECB.

There is plenty of literature explaining why ECB should not be used, starting with https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation, section *Electronic Codebook (ECB)*.

Nevertheless, if like me, someone needed to solve a problem with legacy software and use ECB, this code might be helpful.

## Examples

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

  recoveredPt := decrypt(ct, key)
  fmt.Printf("Recovered plaintext: %s\n", recoveredPt)
}

func main() {
  fmt.Println("AES encryption with ECB and PKCS7 padding")
  example()
}
```

## Padding

Both ECB (Electronic Codebook) and CBC (Cipher Block Chaining) require blocks of fixed size. Encryption with these modes of operation, ECB and CBC, requires to `pad` the plain text to a size multiple of the block size.

The `padding` package exposes simple functions to provide a way to `pad` and `unpad` a given plaintext respectively before encryption and after decryption.

The code examples in the previous sections show encryption patterns with Blowfish and AES in ECB mode. Blowfish encrypts blocks of 8 bytes hence using the padding type described in the https://tools.ietf.org/html/rfc2898 *PKCS #5: Password-Based Cryptography Specification Version 2.0*. Whereas AES requires blocks of 16 bytes (128 bits). The padding type in the second example is based on https://tools.ietf.org/html/rfc2315 *PKCS #7: Cryptographic Message Syntax Version 1.5*.

The only difference between the two specs is that PKCS #5 accommodates only for blocks of 8 bytes. The `padding` package reflects that and exposes two builders, respectively `NewPkcs5Padding()` that embeds a hard-coded value for a block size of 8, while `NewPkcs7Padding(int blockSize)` takes a parameter for the block size. Nothing prevents using `NewPkcs7Padding` with a block size of 8 to work with an encryption scheme working on blocks of 8 bytes, like *Blowfish*.

### Full block of padding

To ensure that there is no ambiguity for the receiver of a message, padding is always performed, even if the message is of an exact multiple block size. This is intentional and comply with the NIST recommendations, * https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf **Appendix A: Padding**, and the following RFCs: 

* https://tools.ietf.org/html/rfc8018
* https://tools.ietf.org/html/rfc1423
* https://tools.ietf.org/html/rfc2898
* https://tools.ietf.org/html/rfc5652

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

The last example is important. Yes, it will add a full block of padding, and this is intentional. This removes the ambiguity for the receiver that expects every message to be padded. 

To illustrate what would happen if some messages are not padded, let's take an example of a message with the last octet with value `01`. As the receiver of this message, should I remove the padding `01`? Or, is the last byte `01`, part of a message that was not padded because it was an exact multiple block size?

If the receiver knows that every message is padded, even if this results in a message padded with a full block of `08`, there is no ambiguity.

Another approach to remove this ambiguity would be to provide a separate indicator that would remove this ambiguity. An example would be to provide a message length indicator.

The implementation in this package relies on padding every message (see method `padding.Pad()`).

## Additional Examples

See the unit tests or the example tests in the respective package directories.

## License

The `ecb` package, is directly modeled after the CBC Golang code (https://golang.org/src/crypto/cipher/cbc.go) released under a BSD license (https://golang.org/LICENSE). To avoid any license conflicts, `crypt2go` is also released under a BSD license.

See the LICENSE file in the repository.



