package main

import (
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

var errInvalidCompressedData = errors.New("invalid compressed data")

type node struct {
	Frequency uint32 `json:"fr,omitempty"`
	Left      *node  `json:"l,omitempty"`
	Right     *node  `json:"r,omitempty"`
	IsLeaf    bool   `json:"ilf,omitempty"`
	Value     byte   `json:"v,omitempty"`
}

// this needs to use a heap for nodes
func buildHuffmanTree(r io.Reader) (*node, error) {
	if r == nil {
		return nil, nil
	}

	byteFrequency := map[byte]uint32{}
	for {
		buff := make([]byte, buffSize)
		n, err := r.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		for i := 0; i < n; i++ {
			byteFrequency[buff[i]]++
		}

		if errors.Is(err, io.EOF) {
			break
		}
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

	return &(*nodes)[0], nil
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

func (n *node) compress(r io.Reader) ([]byte, error) {
	prefixCodeTable := n.buildPrefixCodeTable()

	bits := []bool{}
	for {
		data := make([]byte, buffSize)
		n, err := r.Read(data)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		for i := 0; i < n; i++ {
			code, ok := prefixCodeTable[data[i]]
			if !ok {
				// should never happen
				return nil, errors.New("byte '%c' is not found in prefix code table")
			}

			bits = append(bits, code...)
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	binBytes := make([]byte, int(math.Ceil(float64(len(bits))/8)))
	for idx := range bits {
		currentByte := &binBytes[idx/8]

		if bits[idx] {
			*currentByte += 1 << (7 - idx%8)
		}
	}

	result := make([]byte, 4+len(binBytes))
	binary.BigEndian.PutUint32(result, uint32(len(bits)))
	copy(result[4:], binBytes)

	return result, nil
}

func (root *node) decompress(bits []bool, w io.Writer) error {
	if len(bits) == 0 {
		return nil
	}

	data := make([]byte, 0, buffSize)
	cur := root
	for i := 0; i <= len(bits); i++ {
		if cur == nil {
			return fmt.Errorf("invalid char code: %w", errInvalidCompressedData)
		}

		if cur.IsLeaf {
			data = append(data, cur.Value)
			if len(data) == buffSize {
				_, err := w.Write(data)
				if err != nil {
					return fmt.Errorf("write failed: %w", err)
				}

				data = make([]byte, 0, buffSize)
			}

			cur = root
			if i == len(bits) {
				break
			}

			i--
			continue
		}

		if i == len(bits) {
			// last bit must be a leaf
			return fmt.Errorf("incorrect last character code: %w", errInvalidCompressedData)
		}

		if !bits[i] {
			cur = cur.Left
			continue
		}

		cur = cur.Right
	}

	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}
