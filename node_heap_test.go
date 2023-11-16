package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	h := &nodeHeap{}
	heap.Init(h)
	heap.Push(h, node{Frequency: 10, Value: 1})
	heap.Push(h, node{Frequency: 20, IsLeaf: true, Value: 2})
	heap.Push(h, node{Frequency: 20, Value: 3})
	heap.Push(h, node{Frequency: 15, Value: 5})
	heap.Push(h, node{Frequency: 13, Value: 2})

	pop := heap.Pop(h)
	assert.Equal(t, node{Frequency: 10, Value: 1}, pop)
	assert.Equal(t, h.Len(), 4)

	pop = heap.Pop(h)
	assert.Equal(t, node{Frequency: 13, Value: 2}, pop)
	assert.Equal(t, h.Len(), 3)

	pop = heap.Pop(h)
	assert.Equal(t, node{Frequency: 15, Value: 5}, pop)
	assert.Equal(t, h.Len(), 2)

	pop = heap.Pop(h)
	assert.Equal(t, node{Frequency: 20, IsLeaf: true, Value: 2}, pop)
	assert.Equal(t, h.Len(), 1)

	pop = heap.Pop(h)
	assert.Equal(t, node{Frequency: 20, Value: 3}, pop)
	assert.Equal(t, h.Len(), 0)
}
