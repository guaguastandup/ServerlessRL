package main

import (
	"container/heap"
	"math/rand"
)

var EvictedMemory int64 = 0

var h *IdleContainerHeap = &IdleContainerHeap{}

var totalFrequency int64 = 0

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
	for s.totalMemUsing+3*1024 > s.MEMCapacity {
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
	case "score1": // best
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		score = memory + percentage*100.0
	case "score2":
		// lfu:
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(IntervalCnt[appID]) / float64(totalFrequency)
		score = memory + percentage*80.0 - frequency*200.0
	case "score3":
		// lfu:
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(IntervalCnt[appID]) / float64(totalFrequency)
		score = memory + percentage*100.0 - frequency*100.0
	case "score4":
		// lfu:
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(IntervalCnt[appID]) / float64(totalFrequency)
		score = memory + percentage*50.0 - frequency*300.0
	case "score5":
		// lfu:
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(IntervalCnt[appID]) / float64(totalFrequency)
		score = 1.5*memory + percentage*50.0 - frequency*500.0
	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
