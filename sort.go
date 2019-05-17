package rund

import "container/heap"

// Least Lexicographic Topological Sort.
// graph represents edges, ops represents nodes.
func topoSort(graph map[string][]string, ops map[string]Operator) []string {
	indegree, order := make(map[string]int, len(ops)), make([]string, 0, len(ops))

	for node := range ops {
		indegree[node] = 0
	}
	for _, edges := range graph {
		for _, node := range edges {
			indegree[node]++
		}
	}

	h := &minHeap{}
	heap.Init(h)
	for node, count := range indegree {
		if count == 0 {
			heap.Push(h, node)
		}
	}

	var name string
	for h.Len() > 0 {
		name = heap.Pop(h).(string)
		order = append(order, name)
		for _, node := range graph[name] {
			indegree[node]--
			if indegree[node] == 0 {
				heap.Push(h, node)
			}
		}
	}

	return order
}

type (
	minHeap []string
)

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
