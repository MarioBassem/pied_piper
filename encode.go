package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const buffSize = 16 * 1024

func encode(r io.ReadSeeker) ([]byte, error) {
	// build huffman tree
	tree, err := buildHuffmanTree(r)
	if err != nil {
		return nil, fmt.Errorf("failed to build huffman tree: %w", err)
	}

	treeJson, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("failed to encode tree: %w", err)
	}

	byteCount := uint32(len(treeJson))

	// write count
	// write treeJson
	// use tree to write bytes

	encoding := make([]byte, 4+len(treeJson))
	binary.BigEndian.PutUint32(encoding, byteCount)
	copy(encoding[4:], treeJson)

	_, err = r.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	compressedBytes, err := tree.compress(r)
	if err != nil {
		return nil, err
	}

	encoding = append(encoding, compressedBytes...)

	return encoding, nil
}
