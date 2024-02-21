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
		default:
			container = heap.Pop(h).(*Container)
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

func (s *Server) getScore(appID string, timestamp int64) float64 {
	score := 0.0
	switch policy {
	case "lru":
		break
	case "random":
		score = float64(rand.Intn(10000))
	case "maxmem":
		score = float64(MemoryMap[appID])
	case "maxKeepAlive":
		score = float64(s.currTime - LastIdleTime[appID])
	case "minKeepAlive":
		score = -float64(s.currTime - LastIdleTime[appID])
	case "maxUsage":
		if AppRunningTimeUsage[appID] == 0 {
			score = 0
		} else {
			score = float64(AppRunningMemUsage[appID] / AppRunningTimeUsage[appID])
		}
	case "minUsage":
		if AppRunningTimeUsage[appID] == 0 {
			score = 0
		} else {
			score = -float64(AppRunningMemUsage[appID]) / float64(AppRunningTimeUsage[appID])
		}
	case "maxColdStartRate":
		score = 1.0 - (float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID]))
	case "minColdStartRate":
		score = float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID])
	case "score": // 内存大小 * 到达的平均间隔
		avg := float64(0)
		if IntervalCnt[appID] != 0 {
			avg = float64(IntervalSum[appID] / float64(IntervalCnt[appID]))
		}
		score = float64(MemoryMap[appID]) * avg
	case "score1": // 内存大小 * 到达的平均间隔 * 保持活跃的时间百分比
		avg := float64(0)
		if IntervalCnt[appID] != 0 {
			avg = float64(IntervalSum[appID] / float64(IntervalCnt[appID]))
		}
		score = float64(MemoryMap[appID]) * avg
		keepAliveTime := s.currTime - LastIdleTime[appID]
		percentage := getPercentage(appID, keepAliveTime)
		score = score * percentage
	case "score2": // 内存大小 * 到达的平均间隔 * 保持活跃的时间百分比 / 热启动率, 热启动率越大, 分数越小
		avg := float64(0)
		if IntervalCnt[appID] != 0 {
			avg = float64(IntervalSum[appID] / float64(IntervalCnt[appID]))
		}
		score = float64(MemoryMap[appID]) * avg

		keepAliveTime := s.currTime - LastIdleTime[appID]
		percentage := getPercentage(appID, keepAliveTime)
		warmstart_Rate := float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID])
		if warmstart_Rate == 0 {
			warmstart_Rate = 10000000000.0
		}
		score = score * percentage / warmstart_Rate
	case "score3":
		if IntervalCnt[appID] != 0 {
			score = float64(MemoryMap[appID]) / float64(IntervalCnt[appID]) // cnt大, 分数小
		} else {
			score = float64(MemoryMap[appID])
		}
	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
