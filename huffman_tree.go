package main

import (
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/bits-and-blooms/bitset"
)

var errInvalidCompressedData = errors.New("invalid compressed data")

type node struct {
	Frequency uint32 `json:",omitempty"`
	Left      *node  `json:",omitempty"`
	Right     *node  `json:",omitempty"`
	IsLeaf    bool   `json:",omitempty"`
	Value     byte   `json:",omitempty"`
}

// this needs to use a heap for nodes
func buildHuffmanTree(b []byte) *node {
	if len(b) == 0 {
		return nil
	}

	byteFrequency := map[byte]uint32{}
	for idx := range b {
		byteFrequency[b[idx]]++
	}

	// ensure there are more than one character
	for _, c := range []byte{'a', 'b'} {
		if _, ok := byteFrequency[c]; !ok {
			byteFrequency[c] = 0
		}
	}

	nodes := &nodeHeap{}
	for b, freq := range byteFrequency {
		heap.Push(nodes, node{
			Frequency: freq,
			IsLeaf:    true,
			Value:     b,
		})
	}

	for {
		if nodes.Len() == 1 {
			break
		}

		n1 := heap.Pop(nodes).(node)
		n2 := heap.Pop(nodes).(node)
		heap.Push(nodes, node{
			Left:      &n1,
			Right:     &n2,
			Frequency: n1.Frequency + n2.Frequency,
		})
	}

	return &(*nodes)[0]
}

func (n *node) buildPrefixCodeTable() map[byte][]bool {
	table := map[byte][]bool{}
	representation := []bool{}

	explore(n, representation, table)

	return table
}

// explore explores tree nodes, while assigning binary representation to leaf nodes in the code table
func explore(n *node, representation []bool, table map[byte][]bool) {
	if n == nil {
		// should never happen
		return
	}

	defer func() {
		representation = representation[:len(representation)-1]
	}()

	if n.IsLeaf {
		table[n.Value] = representation
		return
	}

	representation = append(representation, false)
	explore(n.Left, representation, table)
	representation = representation[:len(representation)-1]

	representation = append(representation, true)
	explore(n.Right, representation, table)
	representation = representation[:len(representation)-1]
}

func (n *node) compress(data []byte) ([]byte, error) {
	prefixCodeTable := n.buildPrefixCodeTable()

	bits := []bool{}
	for idx := range data {
		code := prefixCodeTable[data[idx]]
		for _, bit := range code {
			bits = append(bits, bit)
		}
	}

	biteSet := bitset.New(uint(len(bits)))
	for idx := range bits {
		if bits[idx] {
			biteSet.Set(uint(idx))
		}
	}

	bin, err := biteSet.MarshalBinary()
	if err != nil {
		return nil, err
	}
	fmt.Printf("generated bin: %+v", bin)

	result := make([]byte, 4+len(bin))
	copy(result[4:], bin)

	binary.BigEndian.PutUint32(result, uint32(len(bits)))
	return result, nil
}

func (root *node) decompress(bits []bool) ([]byte, error) {
	data := []byte{}
	if len(bits) == 0 {
		return data, nil
	}

	cur := root
	for i := 0; i <= len(bits); i++ {
		if cur == nil {
			return nil, fmt.Errorf("invalid char code: %w", errInvalidCompressedData)
		}

		if cur.IsLeaf {
			data = append(data, cur.Value)
			cur = root
			if i == len(bits) {
				break
			}

			i--
			continue
		}

		if i == len(bits) {
			// last bit must be a leaf
			return nil, fmt.Errorf("incorrect last character code: %w", errInvalidCompressedData)
		}

		if bits[i] == false {
			cur = cur.Left
			continue
		}

		cur = cur.Right
	}

	return data, nil
}
