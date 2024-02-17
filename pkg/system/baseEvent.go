package main

import (
	"container/heap"
)

type event interface {
	getTimestamp() int64
	setTimestamp(int64)
	getHeapIdx() int
	setHeapIdx(int)
	log()
	String() string
}

type baseEvent struct {
	id        int64
	timestamp int64
	heapIdx   int
}

func (e *baseEvent) getTimestamp() int64 {
	return e.timestamp
}

func (e *baseEvent) setTimestamp(t int64) {
	e.timestamp = t
}

func (e *baseEvent) getHeapIdx() int {
	return e.heapIdx
}

func (e *baseEvent) setHeapIdx(i int) {
	e.heapIdx = i
}

type eventQueue []event

func (s *Server) newEventId() int64 {
	s.currEventId++
	return s.currEventId

}
func (s *Server) addEvent(e event) {
	heap.Push(&s.EventQueue, e)
}

func (pq eventQueue) Len() int {
	return len(pq)
}

func (pq eventQueue) Less(i, j int) bool {
	return pq[i].getTimestamp() < pq[j].getTimestamp()
}

func (pq eventQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].setHeapIdx(i)
	pq[j].setHeapIdx(j)
}

func (pq *eventQueue) Push(x interface{}) {
	n := len(*pq)
	e := x.(event)
	e.setHeapIdx(n)
	*pq = append(*pq, e)
}

func (pq *eventQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	e := old[n-1]
	old[n-1] = nil
	e.setHeapIdx(-1)
	*pq = old[0 : n-1]
	return e
}
