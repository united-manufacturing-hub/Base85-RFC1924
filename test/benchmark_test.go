package test

import (
	"base85RFC1924/pkg/base85_rfc1924"
	"base85RFC1924/test/go-rfc1924/base85"
	hailongz "github.com/hailongz/golang/basex"
	"testing"
)

func BenchmarkHailongzBaseX(b *testing.B) {
	encoding, err := hailongz.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~")
	if err != nil {
		b.Fatalf("Error creating encoding: %v", err)
	}
	for i := 0; i < b.N; i++ {
		var decoded []byte
		for _, input := range InputResultMap {
			encoded := encoding.Encode([]byte(input))
			decoded, err = encoding.Decode(encoded)
			if err != nil {
				b.Fatalf("Error decoding: %v", err)
			}
			if string(decoded) != input {
				b.Fatalf("Decoded string does not match input: %v != %v", string(decoded), input)
			}
		}
	}
}

func BenchmarkJamesRuanGoRFC1924(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, input := range InputResultMap {
			encoded := base85.EncodeToString([]byte(input))
			decoded, err := base85.DecodeString(encoded)
			if err != nil {
				b.Fatalf("Error decoding: %v", err)
			}

			if string(decoded) != input {
				b.Fatalf("Decoded string does not match input: %v != %v", string(decoded), input)
			}

		}
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, input := range InputResultMap {
			encoded := base85_rfc1924.EncodeToString([]byte(input))
			decoded, err := base85_rfc1924.DecodeString(encoded)
			if err != nil {
				b.Fatalf("Error decoding: %v", err)
			}

			if string(decoded) != input {
				b.Fatalf("Decoded string does not match input: %v != %v", string(decoded), input)
			}

		}
	}
}
