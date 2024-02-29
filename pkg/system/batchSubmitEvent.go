package main

import (
	"fmt"
	"math"
)

var IntervalSum map[string]float64 = make(map[string]float64)
var IntervalCnt map[string]int = make(map[string]int)

// ***************************** Submit *********************************************
func (s *Server) handleBatchFuncSubmitEvent(e *BatchFunctionSubmitEvent) {
	fmt.Printf("Batch Submit Event: %d day %d minute\n\n", e.day, e.minute)
	if e.minute == 1 {
		initMap()
		ParseMemory(e.day)
		ParseDuration(e.day)
	}
	requests := ParseRequests(e.day, e.minute)
	//增加一份preTime的拷贝, 对PreMinTime修改时, 不会修改preTime
	preMinTime := make(map[string]int64)
	for k, v := range preTime {
		preMinTime[k] = v
	}

	//! batchFunctionSubmit -> functionSubmit
	for _, req := range requests {
		if appHistogram[req.AppID] == nil {
			appHistogram[req.AppID] = &histogram{
				sum:   0,
				array: make([]int, histogramLength),
			}
		}
		if preMinTime[req.AppID] != 0 {
			interval := float64(req.ArrivalTime - preMinTime[req.AppID])
			if interval < 0 {
				interval = 0
			}
			if interval >= 0 {
				interval_min := int(math.Ceil(interval / unit))
				updateHistogram(req.AppID, interval_min)
				IntervalSum[req.AppID] += interval
				IntervalCnt[req.AppID] += 1
			}
		}
		preTime[req.AppID] = req.ArrivalTime
		s.addEvent(&FunctionSubmitEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: req.ArrivalTime,
			},
			function: &Function{
				FuncID:   req.FuncID,
				AppID:    req.AppID,
				FuncType: req.FuncType,
				RunTime:  req.RunTime,
			},
		})
	}
}
