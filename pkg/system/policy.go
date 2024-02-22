package main

import (
	"container/heap"
	"math"
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
	for s.totalMemUsing > s.MEMCapacity {
		if ContainerIdleList.Len() == 0 {
			// panic("No idle container to evict!  " + strconv.Itoa(cnt))
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
	avg := float64(IntervalSum[appID] / float64(IntervalCnt[appID]))
	warmstart_Rate := 1.0 + float64(s.appWarmStartCnt[appID])/float64(s.appRequestCnt[appID])
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
	case "maxmem2":
		score = float64(MemoryMap[appID])*100000000000 + float64(s.appRequestCnt[appID])
	case "maxUsage":
		score = float64(AppRunningMemUsage[appID]) / (float64(AppMemUsage[appID] + 1))
	case "minUsage":
		score = -float64(AppRunningMemUsage[appID]) / (float64(AppMemUsage[appID] + 1))
	case "maxColdStart":
		score = 1.0 - (float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID]))
	case "minColdStart":
		score = float64(s.appWarmStartCnt[appID]) / float64(s.appRequestCnt[appID])
	case "score": // 内存大小 * 到达的平均间隔
		score = float64(MemoryMap[appID]) * avg
	case "score1": // 内存大小 * 到达的平均间隔 * 保持活跃的时间百分比
		score = float64(MemoryMap[appID])
		keepAliveTime := (s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, keepAliveTime)
		score = score * percentage
	case "score2": // 内存大小 * 到达的平均间隔 * 保持活跃的时间百分比 / 热启动率, 热启动率越大, 分数越小
		score = float64(MemoryMap[appID])
		keepAliveTime := (s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, keepAliveTime)
		score = score * percentage / warmstart_Rate
	case "score3":
		score = float64(MemoryMap[appID]) / math.Sqrt((float64(IntervalCnt[appID] + 1))) // cnt大, 分数小
	case "score4":
		interval := s.currTime - LastIdleTime[appID]
		score = -(float64(KeepAliveTimeMap[appID]) - float64(interval)) / float64(KeepAliveTimeMap[appID]) // 剩余时间越大, 优先级越低
	case "score5":
		score = avg * float64(MemoryMap[appID]) * float64(MemoryMap[appID]) / math.Sqrt((float64(IntervalCnt[appID] + 1))) // cnt大, 分数小
	case "score6":
		score = float64(MemoryMap[appID])
		keepAliveTime := (s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, keepAliveTime)
		score = score * percentage * warmstart_Rate
	case "score7": // 目前最佳
		keepAliveTime := (s.currTime - LastIdleTime[appID])
		score = getPercentage(appID, keepAliveTime)
	case "score8":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := math.Pow(float64(MemoryMap[appID]), 1.5)
		score = memory + percentage*100
	case "score9":
		interval := int64(s.currTime - LastIdleTime[appID])
		percentage := getPercentage(appID, interval)
		memory := math.Pow(float64(MemoryMap[appID]), 1.5)
		score = memory + percentage*100
	case "cv":
		interval := s.currTime - LastIdleTime[appID]
		percentage := getPercentage(appID, interval)
		percentage = math.Pow(percentage, 0.5)
		memory := math.Pow(float64(MemoryMap[appID]), 1.5)
		b := math.Pow(float64(IntervalCnt[appID]), 0.5)
		score = memory * percentage / b
		// cv := getCV(appID)
		// if cv == -1 {
		// 	score = memory * percentage
		// } else if cv <= 0.1 {
		// 	score = memory * percentage * 0.5
		// } else if cv <= 0.3 {
		// 	score = memory * percentage * 0.9
		// } else if cv <= 0.5 {
		// 	score = memory * percentage * 1.3
		// } else if cv <= 0.7 {
		// 	score = memory * percentage * 1.7
		// } else {
		// 	score = memory * percentage * 2.1
		// }

	default:
		panic("Unknown policy! " + policy)
	}
	return score
}
