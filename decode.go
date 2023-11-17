package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
)

func decode(r io.Reader, w io.Writer) error {
	// read tree's byte count (n)
	// read next n bytes
	// unmarshal to huffman tree

	root, err := readTree(r)
	if err != nil {
		return err
	}

	// read bin length (m)
	// read next m bits
	// decode from tree

	bits, err := readBits(r)
	if err != nil {
		return err
	}

	if err := root.decompress(bits, w); err != nil {
		return fmt.Errorf("failed to decompress binary data: %w", err)
	}

	return nil
}

func readTree(r io.Reader) (*node, error) {
	count := make([]byte, 4)
	_, err := io.ReadFull(r, count)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree bytes count: %w", err)
	}

	treeBytesCount := binary.BigEndian.Uint32(count)

	treeBytes := make([]byte, treeBytesCount)
	_, err = io.ReadFull(r, treeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree bytes: %w", err)
	}

	root := &node{}
	if err := json.Unmarshal(treeBytes, root); err != nil {
		return nil, err
	}

	return root, nil
}

func readBits(r io.Reader) ([]bool, error) {
	count := make([]byte, 4)
	_, err := io.ReadFull(r, count)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed data size: %w", err)
	}

	binRepBitsCount := binary.BigEndian.Uint32(count)

	bits, err := extractBitsFromBytes(r, int(binRepBitsCount))
	if err != nil {
		return nil, err
	}

	return bits, nil
}

// extractBitsFromBytes extracts the desired number of bits from given data into a bool slice
func extractBitsFromBytes(r io.Reader, bitsCount int) ([]bool, error) {
	bits := make([]bool, 0, bitsCount)

	for {
		data := make([]byte, buffSize)
		n, err := r.Read(data)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to read compressed data: %w", errInvalidCompressedData)
		}

		if int(math.Ceil(float64(bitsCount)/8)) < n+int(math.Ceil(float64(len(bits)/8))) {
			return nil, fmt.Errorf("compressed data size exceeded given size: %w", errInvalidCompressedData)
		}

		data = data[:n]
		for _, b := range data {
			for i := 7; i >= 0 && len(bits) < int(bitsCount); i-- {
				if (1<<i)&b > 0 {
					bits = append(bits, true)
					continue
				}

				bits = append(bits, false)
			}
		}

		if errors.Is(err, io.EOF) {
			if len(bits) < bitsCount {
				return nil, fmt.Errorf("compressed data size is less that given size: %w", errInvalidCompressedData)
			}

			break
		}
	}

	return bits, nil
}
