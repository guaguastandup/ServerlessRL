package main

import (
	"fmt"
)

var AppFuncMap map[string]map[string]int = make(map[string]map[string]int) // appID -> functionID -> start time

// var timeAppUsage map[string]map[int64]int = make(map[string]map[int64]int)       // appID -> time
// var timestampVistedMap map[string]map[int64]int = make(map[string]map[int64]int) // appID -> timestamp
// var timestampMap map[string][]int64 = make(map[string][]int64)                   // appID -> timestamp

func (s *Server) handleFuncStartEvent(e *FunctionStartEvent) {
	e.app.FunctionCnt += 1
	// memory usage
	if AppFuncMap[e.app.AppID] == nil {
		AppFuncMap[e.app.AppID] = make(map[string]int)
	}
	AppFuncMap[e.app.AppID][e.function.FuncID] += 1
	if AppFuncMap[e.app.AppID][e.function.FuncID] == 1 {
		s.totalMemRunning += int64(MemoryFuncMap[e.function.AppID])
	}
	// time usgae
	// if timestampVistedMap[e.app.AppID][e.getTimestamp()] == 0 {
	// 	timestampMap[e.app.AppID] = append(timestampMap[e.app.AppID], e.getTimestamp())
	// 	timestampVistedMap[e.app.AppID][e.getTimestamp()] = 1
	// }
	// timeAppUsage[e.app.AppID][e.getTimestamp()] += 1

	if e.app.FunctionCnt > 0 && e.app.Left == -1 {
		e.app.Left = e.getTimestamp()
	}

	//! functionStart -> functionFinish
	s.addEvent(&FunctionFinishEvent{
		baseEvent: baseEvent{
			id:        s.newEventId(),
			timestamp: e.getTimestamp() + int64(e.function.RunTime),
		},
		function:  e.function,
		app:       e.app,
		container: e.container,
	})
}

func (s *Server) handleFuncFinishEvent(e *FunctionFinishEvent) {
	e.app.FunctionCnt -= 1
	e.app.LastIdleTime = e.getTimestamp() // 记录App最后一次空闲时间

	AppFuncMap[e.app.AppID][e.function.FuncID] -= 1
	if AppFuncMap[e.app.AppID][e.function.FuncID] == 0 {
		s.totalMemRunning -= int64(MemoryFuncMap[e.function.AppID])
	}

	// timeAppUsage[e.app.AppID][e.getTimestamp()] -= 1
	// if timestampVistedMap[e.app.AppID][e.getTimestamp()] == 0 {
	// 	timestampMap[e.app.AppID] = append(timestampMap[e.app.AppID], e.getTimestamp())
	// 	timestampVistedMap[e.app.AppID][e.getTimestamp()] = 1
	// }

	if e.app.FunctionCnt == 0 && e.app.Left != -1 {
		// s.TimeRunningUsage += (e.getTimestamp() - e.app.Left)
		e.app.RunningGain += (e.getTimestamp() - e.app.Left)
		e.app.Left = -1
	}

	//! functionFinish -> appTryFinish
	if e.getTimestamp()+int64(e.app.KeepAliveTime) > e.app.FinishTime { // 产生了新的结束时间
		e.app.FinishTime = e.getTimestamp() + int64(e.app.KeepAliveTime)
		s.addEvent(&AppFinishEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp() + int64(e.app.KeepAliveTime),
			},
			app:       e.app,
			container: e.container,
		})
	}
}

func (s *Server) handleAppInitEvent(e *AppInitEvent) { // 冷启动
	cont := s.NewContainer(e)
	e.app.InitTimeStamp = e.getTimestamp()
	e.app.InitDoneTimeStamp = e.getTimestamp() + int64(e.app.InitTime)
	s.totalMemUsing += int64(e.app.MEMResources)

	// timeAppUsage[e.app.AppID] = make(map[int64]int)
	// timestampMap[e.app.AppID] = make([]int64, 0)
	// timestampVistedMap[e.app.AppID] = make(map[int64]int)

	//! appInit -> functionStart
	s.addEvent(&FunctionStartEvent{
		baseEvent: baseEvent{
			id:        s.newEventId(),
			timestamp: e.getTimestamp() + int64(e.app.InitTime),
		},
		function:  e.function,
		app:       e.app,
		container: cont,
	})
}

func (s *Server) handleAppFinishEvent(e *AppFinishEvent) { // 销毁容器
	if e.app == nil || e.app.FunctionCnt != 0 || e.app.FinishTime != e.getTimestamp() {
		return
	}

	// init := 0
	// left := int64(-1)
	// sum := int64(0)
	// for _, k := range timestampMap[e.app.AppID] {
	// 	v := timeAppUsage[e.app.AppID][k]
	// 	init += v
	// 	if left == -1 && init > 0 {
	// 		left = k
	// 	}
	// 	if init == 0 && left != -1 {
	// 		s.TimeRunningUsage += (k - left)
	// 		sum += (k - left)
	// 		left = -1
	// 	}
	// }

	s.TimeRunningUsage += e.app.RunningGain

	s.MemUsage += int64(e.app.MEMResources) * (e.getTimestamp() - e.app.InitTimeStamp)
	s.TimeUsage += e.getTimestamp() - e.app.InitTimeStamp
	s.totalMemUsing -= int64(e.app.MEMResources)

	s.AppContainerMap[e.app.AppID] = nil
	e.app = nil
}

// ***************************** Submit *********************************************
func (s *Server) handleBatchFuncSubmitEvent(e *BatchFunctionSubmitEvent) {
	fmt.Printf("Batch Submit Event: %d day %d minute\n\n", e.day, e.minute)
	if e.minute == 1 {
		initMap()
		ParseMemory(e.day)
		ParseDuration(e.day)
	}
	requests := ParseRequests(e.day, e.minute)
	//! batchFunctionSubmit -> functionSubmit
	for _, req := range requests {
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

func (s *Server) handleFuncSubmitEvent(e *FunctionSubmitEvent) {
	appID := e.function.AppID
	appMemory := MemoryMap[appID]

	s.appRequestCnt[appID] += 1
	s.totalRequest += 1

	if s.AppContainerMap[appID] != nil { // warm start
		container := s.AppContainerMap[appID]
		app := container.App
		startTime := e.getTimestamp()
		if app.InitDoneTimeStamp <= e.getTimestamp() { // 说明容器已经初始化完成
			s.warmStartCnt++
			s.appWarmStartCnt[appID] += 1
		} else {
			startTime = app.InitDoneTimeStamp
		}
		//! functionSubmit -> functionStart
		s.addEvent(&FunctionStartEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: startTime,
			},
			app:       container.App,
			container: container,
			function:  e.function,
		})
	} else { // cold start
		//! functionSubmit -> appInit
		s.addEvent(&AppInitEvent{ // todo: 此处需要考虑delayed hits的问题
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			},
			function: e.function,
			app: &Application{
				AppID:             appID,
				MEMResources:      appMemory,
				FunctionCnt:       0,
				InitTime:          ColdStartTimeMap[appID],
				InitTimeStamp:     e.getTimestamp(),
				InitDoneTimeStamp: e.getTimestamp() + int64(ColdStartTimeMap[appID]),
				KeepAliveTime:     defaultKeepAliveTime,
				FinishTime:        0,
				Left:              int64(-1),
			},
		})
	}
}
