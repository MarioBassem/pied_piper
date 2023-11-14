package main

type nodeHeap []node

func (n nodeHeap) Len() int {
	return len(n)
}

func (n nodeHeap) Less(i, j int) bool {
	if n[i].Frequency < n[j].Frequency || (n[i].Frequency == n[j].Frequency && n[i].Value < n[j].Value) {
		return true
	}

	return false
}

func (n nodeHeap) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n *nodeHeap) Push(x any) {
	*n = append(*n, x.(node))
}

func (n *nodeHeap) Pop() any {
	old := *n
	l := len(old)

	x := old[l-1]
	*n = old[0 : l-1]
	return x
}

func (n *nodeHeap) Top() any {
	l := len(*n)
	return (*n)[l-1]
}
