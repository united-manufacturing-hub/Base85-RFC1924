package test

import (
	"github.com/united-manufacturing-hub/Base85-RFC1924/pkg/base85_rfc1924"
	"github.com/united-manufacturing-hub/Base85-RFC1924/test/go-rfc1924/base85"
	"testing"
)

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
