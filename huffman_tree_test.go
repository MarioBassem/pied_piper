package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFrequencyTree(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected *node
	}{
		"empty": {
			input: "",
			expected: &node{
				Frequency: 0,
				Left: &node{
					Frequency: 0,
					IsLeaf:    true,
					Value:     'a',
				},
				Right: &node{
					Frequency: 0,
					IsLeaf:    true,
					Value:     'b',
				},
			},
		},
		"same_char": {
			input: "aaaa",
			expected: &node{
				Frequency: 4,
				Left: &node{
					Frequency: 0,
					IsLeaf:    true,
					Value:     'b',
				},
				Right: &node{
					Frequency: 4,
					IsLeaf:    true,
					Value:     'a',
				},
			},
		},
		"mixed_chars": {
			input: "abababaa",
			expected: &node{
				Frequency: 8,
				Left: &node{
					Frequency: 3,
					IsLeaf:    true,
					Value:     'b',
				},
				Right: &node{
					Frequency: 5,
					IsLeaf:    true,
					Value:     'a',
				},
			},
		},
		"same_frequency": {
			input: "ababab",
			expected: &node{
				Frequency: 6,
				Left: &node{
					Frequency: 3,
					IsLeaf:    true,
					Value:     'a',
				},
				Right: &node{
					Frequency: 3,
					IsLeaf:    true,
					Value:     'b',
				},
			},
		},
		"many_mixed_chars": {
			input: "abbcccddddeeeee",
			expected: &node{
				Frequency: 15,
				Left: &node{
					Frequency: 6,
					Left: &node{
						Frequency: 3,
						IsLeaf:    true,
						Value:     'c',
					},
					Right: &node{
						Frequency: 3,
						Left: &node{
							Frequency: 1,
							IsLeaf:    true,
							Value:     'a',
						},
						Right: &node{
							Frequency: 2,
							IsLeaf:    true,
							Value:     'b',
						},
					},
				},
				Right: &node{
					Frequency: 9,
					Left: &node{
						Frequency: 4,
						IsLeaf:    true,
						Value:     'd',
					},
					Right: &node{
						Frequency: 5,
						IsLeaf:    true,
						Value:     'e',
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tc.input))
			got, err := buildHuffmanTree(r)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDecompress(t *testing.T) {
	t.Run("ab_tree", func(t *testing.T) {
		root := &node{Frequency: 6, Left: &node{Frequency: 3, IsLeaf: true, Value: 97}, Right: &node{Frequency: 3, IsLeaf: true, Value: 98}}

		tests := map[string]struct {
			input    []bool
			expected []byte
			hasError bool
		}{
			"nil_input": {
				input:    nil,
				expected: []byte{},
				hasError: false,
			},
			"empty_input": {
				input:    []bool{},
				expected: []byte{},
				hasError: false,
			},
			"valid_data": {
				input:    []bool{false, true, false, false, true},
				expected: []byte("abaab"),
				hasError: false,
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got := bytes.NewBuffer(make([]byte, 0, len(tc.expected)))
				w := io.Writer(got)
				err := root.decompress(tc.input, w)
				if tc.hasError {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got.Bytes())
			})
		}

	})

	t.Run("abc_tree", func(t *testing.T) {
		root := &node{
			Frequency: 10,
			Left: &node{
				Frequency: 5,
				IsLeaf:    true,
				Value:     97,
			},
			Right: &node{
				Frequency: 5,
				Left: &node{
					Frequency: 2,
					IsLeaf:    true,
					Value:     98,
				},
				Right: &node{
					Frequency: 3,
					IsLeaf:    true,
					Value:     99,
				},
			},
		}

		tests := map[string]struct {
			input    []bool
			expected []byte
			hasError bool
		}{
			"nil_input": {
				input:    nil,
				expected: []byte{},
				hasError: false,
			},
			"empty_input": {
				input:    []bool{},
				expected: []byte{},
				hasError: false,
			},
			"valid_data": {
				input:    []bool{false, true, false, true, true},
				expected: []byte("abc"),
				hasError: false,
			},
			"invalid_char": {
				input:    []bool{true},
				hasError: true,
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got := bytes.NewBuffer(make([]byte, 0, len(tc.expected)))
				w := io.Writer(got)
				err := root.decompress(tc.input, w)
				if tc.hasError {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got.Bytes())
			})
		}
	})
}

func TestBuildPrefixCodeTable(t *testing.T) {
	tests := map[string]struct {
		input    *node
		expected map[byte][]bool
	}{
		"nil_tree": {
			input:    nil,
			expected: map[byte][]bool{},
		},
		"two_characters": {
			input: &node{Frequency: 10,
				Left: &node{
					Frequency: 5,
					IsLeaf:    true,
					Value:     'a',
				},
				Right: &node{
					Frequency: 5,
					IsLeaf:    true,
					Value:     'b',
				},
			},
			expected: map[byte][]bool{
				'a': {false},
				'b': {true},
			},
		},
		"three_characters": {
			input: &node{
				Frequency: 10,
				Left: &node{
					Frequency: 5,
					IsLeaf:    true,
					Value:     'a',
				},
				Right: &node{
					Frequency: 5,
					Left: &node{
						Frequency: 3,
						IsLeaf:    true,
						Value:     'b',
					}, Right: &node{
						Frequency: 2,
						IsLeaf:    true,
						Value:     'c',
					},
				},
			},
			expected: map[byte][]bool{
				'a': {false},
				'b': {true, false},
				'c': {true, true},
			},
		},
		"four_characters": {
			input: &node{
				Frequency: 10,
				Left: &node{
					Frequency: 30,
					Left: &node{
						Frequency: 10,
						IsLeaf:    true,
						Value:     'z',
					},
					Right: &node{
						Frequency: 20,
						IsLeaf:    true,
						Value:     'x',
					},
				},
				Right: &node{
					Frequency: 5,
					Left: &node{
						Frequency: 3,
						IsLeaf:    true,
						Value:     'b',
					}, Right: &node{
						Frequency: 2,
						IsLeaf:    true,
						Value:     'c',
					},
				},
			},
			expected: map[byte][]bool{
				'z': {false, false},
				'x': {false, true},
				'b': {true, false},
				'c': {true, true},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.buildPrefixCodeTable()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCompress(t *testing.T) {
	root := &node{
		Frequency: 10,
		Left: &node{
			Frequency: 4,
			IsLeaf:    true,
			Value:     'd',
		},
		Right: &node{
			Frequency: 6,
			Left: &node{
				Frequency: 3,
				IsLeaf:    true,
				Value:     'c',
			}, Right: &node{
				Frequency: 3,
				Left: &node{
					Frequency: 1,
					IsLeaf:    true,
					Value:     'a',
				},
				Right: &node{
					Frequency: 2,
					IsLeaf:    true,
					Value:     'b',
				},
			},
		},
	}

	tests := map[string]struct {
		input    []byte
		expected []byte
		hasError bool
	}{
		"empty": {
			input:    nil,
			expected: []byte{0, 0, 0, 0},
		},
		"one_character": {
			input:    []byte("a"),
			expected: []byte{0, 0, 0, 3, 192},
		},
		"two_characters": {
			input:    []byte("ba"),
			expected: []byte{0, 0, 0, 6, 248},
		},
		"three_characters": {
			input:    []byte("dcb"),
			expected: []byte{0, 0, 0, 6, 92},
		},
		"multiple_characters": {
			input:    []byte("abbcccdddd"),
			expected: []byte{0, 0, 0, 19, 223, 212, 0},
		},
		"inexistent_character": {
			input:    []byte("z"),
			hasError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			got, err := root.compress(r)
			if tc.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
