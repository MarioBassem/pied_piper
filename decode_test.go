package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractBitsFromBytes(t *testing.T) {
	tests := map[string]struct {
		count    int
		data     []byte
		expected []bool
		hasError bool
	}{
		"empty": {
			count:    0,
			data:     nil,
			expected: []bool{},
			hasError: false,
		},
		"more_data": {
			count:    8,
			data:     []byte{'a', 'b'},
			expected: nil,
			hasError: true,
		},
		"less_data": {
			count:    10,
			data:     []byte("a"),
			expected: nil,
			hasError: true,
		},
		"valid_data": {
			count:    9,
			data:     []byte{254, 128},
			expected: []bool{true, true, true, true, true, true, true, false, true},
			hasError: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			got, err := extractBitsFromBytes(r, tc.count)
			if tc.hasError {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.expected, got)
			assert.NoError(t, err)
		})
	}
}

func TestReadBits(t *testing.T) {
	tests := map[string]struct {
		input    []byte
		expected []bool
		hasError bool
	}{
		"empty": {
			input:    nil,
			expected: nil,
			hasError: true,
		},
		"invalid_data": {
			input:    []byte{'a', 'b', 'c', 'd', 'e'},
			expected: nil,
			hasError: true,
		},
		"valid_data": {
			input:    []byte{0, 0, 0, 1, 128},
			expected: []bool{true},
			hasError: false,
		},
		"more_binary_data": {
			input:    []byte{0, 0, 0, 1, 192, 1},
			expected: nil,
			hasError: true,
		},
		"less_binary_data": {
			input:    []byte{0, 0, 0, 10, 192},
			expected: nil,
			hasError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			got, err := readBits(r)
			if tc.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestReadTree(t *testing.T) {
	validData := []byte("{\"fr\":6,\"l\":{\"fr\":3,\"ilf\":true,\"v\":97},\"r\":{\"fr\":3,\"ilf\":true,\"v\":98}}")
	tests := map[string]struct {
		input    []byte
		expected *node
		hasError bool
		count    int
	}{
		"empty": {
			input:    nil,
			hasError: true,
		},
		"valid_data": {
			input:    append([]byte{0, 0, 0, byte(len(validData))}, validData...),
			hasError: false,
			expected: &node{Frequency: 6, Left: &node{Frequency: 3, IsLeaf: true, Value: 97}, Right: &node{Frequency: 3, IsLeaf: true, Value: 98}},
			count:    4 + len(validData),
		},
		"bytes_number_less_than_expected": {
			input:    []byte{0, 0, 0, 1, 'a', 'b', 'c'},
			hasError: true,
		},
		"bytes_number_more_than_expected": {
			input:    []byte{0, 0, 0, 10, 'a', 'b'},
			hasError: true,
		},
		"invalid_json_tree": {
			input:    append([]byte{0, 0, 0, 7}, []byte("{hello}")...),
			hasError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			got, err := readTree(r)
			if tc.hasError {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.expected, got)
			assert.NoError(t, err)
		})
	}
}
