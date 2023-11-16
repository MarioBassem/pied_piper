package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

func decode(data []byte) ([]byte, error) {
	// read tree's byte count (n)
	// read next n bytes
	// unmarshal to huffman tree
	treeBytesCount := binary.BigEndian.Uint32(data[:3])
	data = data[4:]

	treeBytes := data[:treeBytesCount]
	root := &node{}
	if err := json.Unmarshal(treeBytes, root); err != nil {
		return nil, err
	}
	root, n, err := readTree(data)
	if err != nil {
		return nil, err
	}

	// read bin length (m)
	// read next m bits
	// decode from tree
	if len(data) < n {
		// this should never happen
		return nil, errInvalidCompressedData
	}
	data = data[n:]

	bits, err := readBits(data)
	if err != nil {
		return nil, err
	}

	decompressedBytes, err := root.decompress(bits)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress binary data: %w", err)
	}

	return decompressedBytes, nil
}

func readTree(data []byte) (*node, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("failed to read tree bytes count: %w", errInvalidCompressedData)
	}

	treeBytesCount := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	if len(data) < int(treeBytesCount) {
		return nil, 0, fmt.Errorf("failed to read tree bytes: %w", errInvalidCompressedData)
	}

	treeBytes := data[:treeBytesCount]
	root := &node{}
	if err := json.Unmarshal(treeBytes, root); err != nil {
		return nil, 0, err
	}

	n := 4 + int(treeBytesCount)
	return root, n, nil
}

func readBits(data []byte) ([]bool, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("failed to read compressed data bits count: %w", errInvalidCompressedData)
	}

	binRepBitsCount := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	bits, err := extractBitsFromBytes(int(binRepBitsCount), data)
	if err != nil {
		return nil, err
	}

	return bits, nil
}

// extractBitsFromBytes extracts the desired number of bits from given data into a bool slice
func extractBitsFromBytes(bitsCount int, data []byte) ([]bool, error) {
	if len(data) != int(math.Ceil(float64(bitsCount)/8)) {
		return nil, fmt.Errorf("failed to read compressed bits: %w", errInvalidCompressedData)
	}

	bits := make([]bool, 0, bitsCount)

	for _, b := range data {
		for i := 7; i >= 0 && len(bits) < int(bitsCount); i-- {
			if (1<<i)&b > 0 {
				bits = append(bits, true)
				continue
			}

			bits = append(bits, false)
		}
	}

	return bits, nil
}
