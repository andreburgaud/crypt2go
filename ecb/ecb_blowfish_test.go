// Copyright 2016 Andre Burgaud. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ecb

import (
	"bytes"
	"testing"

	//nolint:staticcheck // SA1019 ignore this
	"golang.org/x/crypto/blowfish"
)

var ecbBlowfishTests = []ecbTest{

	// NIST SP 800-38A pp 24-27
	{
		"ECB-Blowfish128",
		commonKey128,
		commonInput,
		[]byte{
			0x05, 0xac, 0x65, 0x4c, 0x24, 0x4c, 0x58, 0x65, 0x7b, 0x82, 0x8d, 0xcf, 0xe6, 0xf8, 0x38, 0x9e,
			0x1d, 0xfa, 0x4a, 0x38, 0x17, 0xb1, 0x3e, 0xd3, 0x1a, 0x5d, 0xe0, 0xd6, 0x0f, 0x53, 0xac, 0x22,
			0x4b, 0x46, 0x34, 0xdd, 0xa3, 0xdf, 0xb3, 0xc9, 0x75, 0x70, 0x3f, 0xef, 0xa9, 0xe6, 0x91, 0x56,
			0x0e, 0xa0, 0x10, 0x1b, 0x30, 0x60, 0xb1, 0x16, 0x8a, 0x8c, 0x4a, 0x9f, 0xf6, 0x88, 0x48, 0x16,
		},
	},

	{
		"ECB-Blowfish192",
		commonKey192,
		commonInput,
		[]byte{
			0xda, 0x42, 0x91, 0x7c, 0x1f, 0x7a, 0x20, 0x20, 0xb5, 0x07, 0x08, 0x91, 0x60, 0xef, 0xf0, 0xf1,
			0x74, 0xa0, 0x62, 0x6c, 0xad, 0x63, 0x94, 0xe6, 0x0c, 0xdc, 0x58, 0xf3, 0x12, 0x9e, 0xf5, 0x53,
			0xa3, 0x77, 0x30, 0x66, 0xf9, 0x1d, 0x2c, 0x28, 0x69, 0x92, 0x43, 0xe3, 0xfc, 0xbf, 0x81, 0xdc,
			0x8d, 0x1e, 0x4c, 0x81, 0xe3, 0xab, 0x64, 0xda, 0x6c, 0x71, 0xfe, 0xc5, 0xfd, 0xd4, 0xfb, 0x4c,
		},
	},

	{
		"ECB-Blowfish256",
		commonKey256,
		commonInput,
		[]byte{
			0xd7, 0xf8, 0xbe, 0x57, 0xaa, 0xe1, 0xc6, 0x36, 0xf5, 0x03, 0xae, 0x7f, 0xf6, 0xcb, 0x1d, 0x25,
			0x7c, 0xf7, 0xc9, 0x56, 0xf3, 0x69, 0x5f, 0x2f, 0x08, 0x09, 0x57, 0xa7, 0x44, 0x86, 0x95, 0x7c,
			0x1b, 0xa0, 0x57, 0x64, 0xe8, 0x10, 0x2a, 0xb8, 0x45, 0x24, 0x40, 0x40, 0xd1, 0x7c, 0x97, 0xa0,
			0xf6, 0x5a, 0x48, 0x58, 0x47, 0x82, 0x4d, 0xe4, 0xae, 0x78, 0x7a, 0xc3, 0xe0, 0xa1, 0x4e, 0x48,
		},
	},
}

func TestECBEncrypterBlowfish(t *testing.T) {
	for _, tc := range ecbBlowfishTests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := blowfish.NewCipher(tc.key)
			if err != nil {
				t.Errorf("%s: NewCipher(%d bytes) = %s", tc.name, len(tc.key), err)
			}
			encrypter := NewECBEncrypter(c)
			data := make([]byte, len(tc.in))
			copy(data, tc.in)
			encrypter.CryptBlocks(data, data)
			if !bytes.Equal(tc.out, data) {
				t.Errorf("%s: ECBEncrypter\nhave %x\nwant %x", tc.name, data, tc.out)
			}
		})
	}
}

func TestECBDecrypterBlowfish(t *testing.T) {
	for _, tc := range ecbBlowfishTests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := blowfish.NewCipher(tc.key)
			if err != nil {
				t.Fatalf("%s: NewCipher(%d bytes) = %s", tc.name, len(tc.key), err)
			}
			decrypter := NewECBDecrypter(c)
			data := make([]byte, len(tc.out))
			copy(data, tc.out)
			decrypter.CryptBlocks(data, data)
			if !bytes.Equal(tc.in, data) {
				t.Errorf("%s: ECBDecrypter\nhave %x\nwant %x", tc.name, data, tc.in)
			}
		})
	}
}
