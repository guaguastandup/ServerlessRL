package main

import "container/heap"

type IdleContainerHeap []*Container

var m map[int64]*Container = make(map[int64]*Container)

// heap
func (h IdleContainerHeap) Len() int           { return len(h) }
func (h IdleContainerHeap) Less(i, j int) bool { return h[i].App.Score > h[j].App.Score } // 比较函数

// Push 向堆中添加元素
func (h *IdleContainerHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*Container)
	item.Index = n // 设置新元素的索引
	m[item.ID] = item
	*h = append(*h, item)
}

func (h IdleContainerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i // 更新交换后的索引
	h[j].Index = j
}

// Pop 从堆中移除并返回最小元素
func (h *IdleContainerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	ele := old[n-1]
	old[n-1] = nil
	ele.Index = -1 // 从堆中移除
	*h = old[0 : n-1]
	delete(m, ele.ID)
	return ele
}

func (h *IdleContainerHeap) RemoveByID(id int64) {
	container, ok := m[id]
	if !ok {
		return // 元素不存在
	}
	index := container.Index
	if index < 0 || index >= h.Len() {
		return
	}
	h.Swap(index, h.Len()-1)
	*h = (*h)[:h.Len()-1]
	if index < h.Len() {
		heap.Fix(h, index)
	}
	delete(m, id)
}

func (h *IdleContainerHeap) UpdateScore(id int64, newScore float64) {
	// Step 1: Locate the element in the heap
	container, ok := m[id]
	if !ok {
		return // Element not found
	}
	// Step 2: Update the key value
	container.App.Score = newScore
	// Step 3: Re-sort the heap
	heap.Fix(h, container.Index)
}
