package main

import (
	"container/heap"

	"github.com/bits-and-blooms/bitset"
)

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
	return nil
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

	binary, err := biteSet.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return binary, nil
}
