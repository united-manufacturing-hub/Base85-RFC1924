package test

import (
	"bytes"
	"github.com/united-manufacturing-hub/Base85-RFC1924/pkg/base85_rfc1924"
	"github.com/united-manufacturing-hub/Base85-RFC1924/test/go-rfc1924/base85"
	"testing"
)

func TestJamesRuanGoRFC1924(t *testing.T) {

	for _, input := range InputResultMap {
		encoded := base85.EncodeToString([]byte(input))
		t.Logf("Encoded %v to %v [%d]", input, encoded, len(encoded))

		decoded, err := base85.DecodeString(encoded)

		if err != nil {
			t.Fatalf("Error decoding: %v", err)
		}

		if string(decoded) != input {
			t.Fatalf("Decoded string does not match input: %v != %v", string(decoded), input)
		}

	}
}

func TestOptimizedB85(t *testing.T) {

	for _, input := range InputResultMap {
		encoded := base85_rfc1924.EncodeToString([]byte(input))
		t.Logf("Encoded %v to %v [%d]", input, encoded, len(encoded))

		decoded, err := base85_rfc1924.DecodeString(encoded)

		if err != nil {
			t.Fatalf("Error decoding: %v", err)
		}

		if string(decoded) != input {
			t.Fatalf("Decoded string does not match input: %v != %v", string(decoded), input)
		}

	}
}

func TestOptimizedAgainstBase(t *testing.T) {
	for _, input := range InputResultMap {
		encodedA := base85_rfc1924.EncodeToString([]byte(input))
		encodedB := base85.EncodeToString([]byte(input))

		decodedA, err := base85_rfc1924.DecodeString(encodedA)
		if err != nil {
			t.Fatalf("Error decoding: %v", err)
		}

		decodedB, err := base85.DecodeString(encodedB)

		if err != nil {
			t.Fatalf("Error decoding: %v", err)
		}

		if !bytes.Equal(decodedA, decodedB) {
			t.Fatalf("Decoded string does not match input: %v != %v", string(decodedA), string(decodedB))
		}

		if encodedA != encodedB {
			t.Fatalf("Encoded string does not match input: %v != %v", encodedA, encodedB)
		}
	}
}
