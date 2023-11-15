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
