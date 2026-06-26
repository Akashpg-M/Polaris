package graph

// import "container/heap"

// pqItem represents a node in the priority queue
type pqItem struct {
	nodeID int64
	cost   float64 // The accumulated traversal cost to reach this node
	index  int     // The index of the item in the heap (needed for internal heap operations)
}

// priorityQueue implements heap.Interface and holds pqItems
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

// Less ensures we pop the node with the lowest cost (Min-Heap)
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	item.index = -1 // For safety
	*pq = old[0 : n-1]
	return item
}