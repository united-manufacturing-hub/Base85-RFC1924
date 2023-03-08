package test

import (
	"bytes"
	"fmt"
	"github.com/united-manufacturing-hub/Base85-RFC1924/test/go-rfc1924/base85"
	"testing"
)

func EncodeDecodeJamesRuanRFC1924(input []byte) {
	encoded := base85.EncodeToString(input)
	decoded, err := base85.DecodeString(encoded)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(decoded, input) {
		panic(fmt.Sprintf("Decoded string does not match input: %v != %v", string(decoded), input))
	}
}

func FuzzEncodeDecode(f *testing.F) {
	for _, input := range InputResultMap {
		f.Add([]byte(input))
	}
	f.Fuzz(func(t *testing.T, input []byte) {
		EncodeDecodeJamesRuanRFC1924(input)
	})
}
