package main

import (
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"
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
		return
	}

	if n.IsLeaf {
		cp := make([]bool, len(representation))
		copy(cp, representation)
		table[n.Value] = cp
		return
	}

	explore(n.Left, append(representation, false), table)
	explore(n.Right, append(representation, true), table)
}

func (n *node) compress(data []byte) ([]byte, error) {
	prefixCodeTable := n.buildPrefixCodeTable()

	bits := []bool{}
	for idx := range data {
		code, ok := prefixCodeTable[data[idx]]
		if !ok {
			// should never happen
			return nil, errors.New("byte '%c' is not found in prefix code table")
		}

		for _, bit := range code {
			bits = append(bits, bit)
		}
	}

	binBytes := []byte{}
	currentByte := byte(0)
	for idx := range bits {
		if idx > 0 && idx%8 == 0 {
			binBytes = append(binBytes, currentByte)
			currentByte = 0
		}

		if bits[idx] {
			currentByte += 1 << (7 - idx%8)
		}
	}

	if len(bits)%8 != 0 {
		binBytes = append(binBytes, currentByte)
	}

	result := make([]byte, 4+len(binBytes))
	binary.BigEndian.PutUint32(result, uint32(len(bits)))
	copy(result[4:], binBytes)

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
