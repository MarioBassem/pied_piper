package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
