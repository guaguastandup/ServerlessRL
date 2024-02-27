package main

import (
	"container/heap"
	"math"
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
	for s.totalMemUsing+1*1024 > s.MEMCapacity {
		if ContainerIdleList.Len() == 0 {
			break
		}
		var container *Container
		switch policy {
		case "lru":
			container = FrontElement()
		case "mru":
			container = BackElement()
		case "ideal":
			container = FrontElement()
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
	case "ideal":
		break
	case "lfu":
		score = -float64(s.appRequestCnt[appID])
	case "random":
		score = float64(rand.Intn(10000))
	case "maxmem":
		score = float64(MemoryMap[appID])
	case "score1":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(frequencyMap[appID]) / float64(totalFrequency)
		score = 1.5*memory + percentage*50.0 - frequency*500.0
	case "score2":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval) * 100.0
		memory := float64(MemoryMap[appID])
		frequency := float64(frequencyMap[appID]) / float64(totalFrequency) * 100.0
		if memory <= 25 {
			score = 2.0*memory + percentage*1.5 - frequency*2.0
		} else if memory <= 50 {
			score = 1.5*memory + percentage*1.25 - frequency*3.0
		} else if memory <= 100 {
			score = 1.0*memory + percentage*1.0 - frequency*4.0
		} else if memory <= 150 {
			score = 1.0*memory + percentage*0.75 - frequency*5.0
		} else {
			score = 1.0*memory + percentage*0.5 - frequency*6.0
		}
	case "score3":
		score = -float64(ColdStartTimeMap[appID])*float64(MemoryMap[appID]) - float64(frequencyMap[appID])*100.0
	case "score4":
		// conbime score1 and score3
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := float64(MemoryMap[appID])
		frequency := float64(frequencyMap[appID]) / float64(totalFrequency)
		if memory <= 25 {
			score = 2.0*memory + percentage*1.5 - frequency*2.0
		} else if memory <= 50 {
			score = 1.5*memory + percentage*1.25 - frequency*3.0
		} else if memory <= 100 {
			score = 1.0*memory + percentage*1.0 - frequency*4.0
		} else if memory <= 150 {
			score = 1.0*memory + percentage*0.75 - frequency*5.0
		} else {
			score = 1.0*memory + percentage*0.5 - frequency*6.0
		}
		memcost := float64(ColdStartTimeMap[appID]) * float64(MemoryMap[appID])
		score -= math.Pow(memcost, 0.25) / 100.0
	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
