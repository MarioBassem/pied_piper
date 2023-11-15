package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func encode(b []byte) ([]byte, error) {
	// build huffman tree
	tree := buildHuffmanTree(b)

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

	compressedBytes, err := tree.compress(b)
	if err != nil {
		return nil, err
	}

	encoding = append(encoding, compressedBytes...)

	return encoding, nil
}
