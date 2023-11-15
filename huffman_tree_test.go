package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFrequencyTree(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected *node
	}{
		"empty": {
			input:    "",
			expected: nil,
		},
		"same_char": {
			input: "aaaa",
			expected: &node{
				Frequency: 4,
				Left:      nil,
				Right:     nil,
				IsLeaf:    true,
				Value:     'a',
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
					Right: &node{
						Frequency: 3,
						IsLeaf:    true,
						Value:     'c',
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
			got := buildHuffmanTree([]byte(tc.input))
			assert.Equal(t, tc.expected, got)

			jsonTree, err := json.Marshal(got)
			assert.NoError(t, err)
			fmt.Printf("json tree: %s", jsonTree)
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
				got, err := root.decompress(tc.input)
				if tc.hasError {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
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
				got, err := root.decompress(tc.input)
				if tc.hasError {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			})
		}
	})
}

func printTree(t *node) {
	if t == nil {
		fmt.Printf("\n")
		return
	}

	fmt.Printf("%+v\n", t)
	printTree(t.Left)
	printTree(t.Right)
}
