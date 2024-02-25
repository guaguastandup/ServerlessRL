package main

import (
	"container/heap"
	"math/rand"
)

var EvictedMemory int64 = 0

var h *IdleContainerHeap = &IdleContainerHeap{}

func (s *Server) handleEvictEvent(e *baseEvent) {
	cnt := 0
	for cont := ContainerIdleList.Front(); cont != nil; cont = cont.Next() {
		container := cont.Value.(*Container)
		score := s.getScore(container.App.AppID, e.getTimestamp())
		if score != container.App.Score {
			h.UpdateScore(container.ID, score)
		}
		cnt += 1
	}
	for s.totalMemUsing+2*1024 > s.MEMCapacity {
		if ContainerIdleList.Len() == 0 {
			break
		}
		var container *Container
		switch policy {
		case "lru":
			container = FrontElement()
		case "mru":
			container = BackElement()
		default:
			container = heap.Pop(h).(*Container)
		}
		cnt -= 1
		RemoveIdleContainer(container)
		EvictedMemory += int64(container.App.MEMResources)
		container.App.FinishTime = e.getTimestamp()
		container.App.PreWarmTime = 0
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

func (s *Server) getScore(appID string, timestamp int64) float64 {
	score := 0.0
	switch policy {
	case "lru":
		break
	case "mru":
		break
	case "lfu":
		score = -float64(s.appRequestCnt[appID])
	case "random":
		score = float64(rand.Intn(10000))
	case "maxmem":
		score = float64(MemoryMap[appID])
	case "score0":
		avg := float64(IntervalSum[appID] / float64(IntervalCnt[appID]))
		score = float64(MemoryMap[appID]) * avg
	case "score1":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		// memory := math.Pow(float64(MemoryMap[appID]), 1.0)
		score = memory + percentage*100
	case "score2":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := float64(interval) / float64(KeepAliveTimeMap[appID])
		memory := float64(MemoryMap[appID])
		// memory := math.Pow(float64(MemoryMap[appID]), 1.0)
		score = memory + percentage*100
	case "score3":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		// memory := math.Pow(float64(MemoryMap[appID]), 1.0)
		score = memory + percentage*200
	case "score4":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		// memory := math.Pow(float64(MemoryMap[appID]), 1.0)
		score = memory + percentage*1000
	case "score5":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		score = percentage
	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
