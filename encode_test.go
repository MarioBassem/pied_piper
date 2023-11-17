package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-={}[]:'\";/?.>,<\\|`~")

func randSeq(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func TestEncode(t *testing.T) {
	tests := map[string]struct {
		input          []byte
		compressedData []byte
		hasError       bool
	}{
		"empty": {
			input:          nil,
			compressedData: []byte{0, 0, 0, 0},
		},
		"single_character": {
			input:          []byte("a"),
			compressedData: []byte{0, 0, 0, 1, 128},
		},
		"multiple_same_character": {
			input:          []byte("aaaa"),
			compressedData: []byte{0, 0, 0, 4, 240},
		},
		"multiple_differenct_characters": {
			input:          []byte("abca"),
			compressedData: []byte{0, 0, 0, 6, 88},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			tree, err := buildHuffmanTree(r)
			assert.NoError(t, err)

			jsonTree, err := json.Marshal(tree)
			assert.NoError(t, err)
			cntBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(cntBytes, uint32(len(jsonTree)))

			want := append(cntBytes, jsonTree...)
			want = append(want, tc.compressedData...)

			r = bytes.NewReader(tc.input)
			got, err := encode(r)
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	for i := 1; i <= 10000000; i *= 10 {
		str := randSeq(i)
		r := bytes.NewReader(str)

		compressed, err := encode(r)
		assert.NoError(t, err)
		fmt.Printf("compressed data length: %d\n", len(compressed))

		got := bytes.NewBuffer(make([]byte, 0, len(str)))
		w := io.Writer(got)

		r2 := bytes.NewReader(compressed)
		err = decode(r2, w)
		assert.NoError(t, err)

		assert.Equal(t, str, got.Bytes())
	}
}
