package main

import (
	"container/heap"
	"math/rand"
)

var EvictedMemory int64 = 0

var h *IdleContainerHeap = &IdleContainerHeap{}

func (s *Server) handleEvictEvent(e *baseEvent) {
	for s.totalMemUsing > s.MEMCapacity {
		if ContainerIdleList.Len() == 0 {
			panic("No idle container to evict!")
		}
		var container *Container
		switch policy {
		case "lru":
			container = FrontElement()
		case "random":
			// index := rand.Intn(ContainerIdleList.Len())
			// container = GetElementByIndex(index)
			container = heap.Pop(h).(*Container)
		case "maxmem":
			// maxMem := 0
			// for ele := ContainerIdleList.Front(); ele != nil; ele = ele.Next() {
			// 	if ele.Value.(*Container).App.MEMResources >= maxMem {
			// 		maxMem = ele.Value.(*Container).App.MEMResources
			// 		container = ele.Value.(*Container)
			// 	}
			// }
			container = heap.Pop(h).(*Container)
		case "maxKeepAlive":
			// maxKeepAlive := 0
			// for ele := ContainerIdleList.Front(); ele != nil; ele = ele.Next() {
			// 	if ele.Value.(*Container).App.KeepAliveTime >= maxKeepAlive {
			// 		maxKeepAlive = ele.Value.(*Container).App.KeepAliveTime
			// 		container = ele.Value.(*Container)
			// 	}
			// }
			container = heap.Pop(h).(*Container)
		case "minUsage":
			// minUsage := int64(1e18)
			// for ele := ContainerIdleList.Front(); ele != nil; ele = ele.Next() {
			// 	appID := ele.Value.(*Container).App.AppID
			// 	if AppRunningMemUsage[appID] < minUsage {
			// 		minUsage = AppRunningMemUsage[appID]
			// 		container = ele.Value.(*Container)
			// 	}
			// }
			container = heap.Pop(h).(*Container)
		case "maxColdStartRate":
			// maxColdStart := float64(0.0)
			// for ele := ContainerIdleList.Front(); ele != nil; ele = ele.Next() {
			// 	appID := ele.Value.(*Container).App.AppID
			// 	if float64(s.appWarmStartCnt[appID])/float64(s.appRequestCnt[appID]) >= maxColdStart {
			// 		maxColdStart = float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID])
			// 		container = ele.Value.(*Container)
			// 	}
			// }
			container = heap.Pop(h).(*Container)
		default:
			panic("Unknown policy! " + policy + container.App.AppID)
		}
		RemoveIdleContainer(container)
		EvictedMemory += int64(container.App.MEMResources)
		container.App.FinishTime = e.getTimestamp()
		s.handleAppFinishEvent(&AppFinishEvent{ // 即刻执行
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			},
			app:       container.App,
			container: container,
		})
	}
}

func (s *Server) getScore(appID string) float64 {
	score := 0.0
	switch policy {
	case "lru":
		break
	case "random":
		score = float64(rand.Intn(10000))
	case "maxmem":
		score = float64(MemoryMap[appID])
	case "maxKeepAlive":
		score = float64(defaultKeepAliveTime)
	case "minUsage":
		score = -float64(AppRunningMemUsage[appID]) / float64(AppRunningTimeUsage[appID])
	case "maxColdStartRate":
		score = 1.0 - (float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID]))
	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
