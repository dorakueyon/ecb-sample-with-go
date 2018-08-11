package tool

import (
	"bytes"
	_ "crypto/des"
	"golang.org/x/crypto/blowfish"
	"testing"
)

func TestEncAndDec(t *testing.T) {
	tests := []struct {
		key   []byte
		input []byte
	}{
		{
			[]byte("priv-key"),
			[]byte("plaintxt"),
		},

		{
			[]byte("priv-key"),
			[]byte("plaintxtplaintxt"),
		},
		{
			// 2byte charater test
			[]byte("priv-key"),
			[]byte("これはテストです。文字数は１６字"),
		},

		{
			[]byte("GN3833WYMjDj4qi7"),
			[]byte{0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08},
		},
	}

	for _, tt := range tests {
		block, err := blowfish.NewCipher(tt.key)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		crypted := make([]byte, len(tt.input))
		ecb := NewECBEncrypter(block)
		ecb.CryptBlocks(crypted, tt.input)

		ecbd := NewECBDecrypter(block)

		decrypted := make([]byte, len(tt.input))
		ecbd.CryptBlocks(decrypted, crypted)

		if ok := bytes.Equal(decrypted, tt.input); !ok {
			t.Fatalf("dec failed. got=:%x,want=%x", decrypted, tt.input)
		}
	}
}

func TestPadding(t *testing.T) {
	tests := []struct {
		input       []byte
		blockSize   int
		expectedVal []byte
	}{
		{[]byte{0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18},
			8,
			[]byte{0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08},
		},
		{[]byte{0x83, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18},
			7,
			[]byte{0x83, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x53, 0x8b, 0x77, 0x59, 0x83, 0x4d, 0x34, 0x18, 0x05, 0x05, 0x05, 0x05, 0x05},
		},
		{[]byte("plaintext"),
			8,
			[]byte{0x70, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07},
		},
	}
	for _, tt := range tests {
		v := PKCS5Padding(tt.input, tt.blockSize)
		if ok := bytes.Equal(v, tt.expectedVal); !ok {
			t.Errorf("Padded val is wrong. got=%x, want=%x, %v", v, tt.expectedVal, v)
		}
	}
}
