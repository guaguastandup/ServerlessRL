package main

import (
	"fmt"
	"math"
)

var AppFuncMap map[string]map[string]int = make(map[string]map[string]int)  // appID -> functionID -> start time
var AppLeft map[string]map[string]int64 = make(map[string]map[string]int64) // appID -> left time
var AppMemUsage map[string]int64 = make(map[string]int64)                   // appID -> mem usage
var AppRunningMemUsage map[string]int64 = make(map[string]int64)            // appID -> running mem usage
var AppTimeUsage map[string]int64 = make(map[string]int64)                  // appID -> time usage
var AppRunningTimeUsage map[string]int64 = make(map[string]int64)           // appID -> running time usage

var ContainerIdleList []*Container
var ContainerIdleMap map[*Container]bool = make(map[*Container]bool) // 一开始全都是空闲的

var EvictedMemory int64 = 0

type histogram struct {
	sum   int
	array []int
}

var appHistogram map[string]*histogram = make(map[string]*histogram)

func (s *Server) handleEvictEvent(e *baseEvent) {
	for s.totalMemUsing > s.MEMCapacity {
		// 删除第一个空闲的容器
		if len(ContainerIdleList) == 0 {
			panic("No idle container to evict!")
		}
		container := ContainerIdleList[0]
		ContainerIdleList = ContainerIdleList[1:]
		if !ContainerIdleMap[container] {
			panic("impossible")
		}
		ContainerIdleMap[container] = false
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

func (s *Server) handleFuncStartEvent(e *FunctionStartEvent) {
	if ContainerIdleMap[e.container] {
		ContainerIdleMap[e.container] = false
		RemoveIdleContainer(e.container)
	}
	AppFuncMap[e.app.AppID][e.function.FuncID] += 1
	if AppFuncMap[e.app.AppID][e.function.FuncID] == 1 {
		s.totalMemRunning += int64(MemoryFuncMap[e.function.AppID])
	}
	if AppFuncMap[e.app.AppID][e.function.FuncID] > 0 && AppLeft[e.app.AppID][e.function.FuncID] == 0 {
		AppLeft[e.app.AppID][e.function.FuncID] = e.getTimestamp()
	}
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
	AppFuncMap[e.app.AppID][e.function.FuncID] -= 1
	e.app.LastIdleTime = e.getTimestamp() // 记录App最后一次空闲时间
	if AppFuncMap[e.app.AppID][e.function.FuncID] == 0 {
		s.totalMemRunning -= int64(MemoryFuncMap[e.function.AppID])
	}
	if AppFuncMap[e.app.AppID][e.function.FuncID] == 0 && AppLeft[e.app.AppID][e.function.FuncID] != 0 {
		e.app.MemRunningGain += (e.getTimestamp() - AppLeft[e.app.AppID][e.function.FuncID]) * int64(MemoryFuncMap[e.function.AppID])
		AppLeft[e.app.AppID][e.function.FuncID] = 0
	}
	if e.app.FunctionCnt == 0 && e.app.Left != -1 {
		e.app.RunningGain += (e.getTimestamp() - e.app.Left)
		e.app.Left = -1
	}
	if e.app.FunctionCnt == 0 {
		if !ContainerIdleMap[e.container] {
			ContainerIdleMap[e.container] = true
			ContainerIdleList = append(ContainerIdleList, e.container)
		}
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
	if appHistogram[e.app.AppID] == nil {
		appHistogram[e.app.AppID] = &histogram{
			sum:   0,
			array: make([]int, 360),
		}
	}
	flag := 0
	if s.AppContainerMap[e.app.AppID] == nil {
		flag = 1
		_ = s.NewContainer(e)
		e.app.InitTimeStamp = e.getTimestamp()
		e.app.InitDoneTimeStamp = e.getTimestamp() + int64(e.app.InitTime)
		s.totalMemUsing += int64(e.app.MEMResources)
		if s.totalMemUsing > s.MEMCapacity { // Memory Overload
			s.handleEvictEvent(&baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			})
		}
		AppFuncMap[e.app.AppID] = make(map[string]int)
		AppLeft[e.app.AppID] = make(map[string]int64)
	}
	if flag == 1 && e.function == nil { // 预热
		e.app.FinishTime = e.getTimestamp() + int64(e.app.InitTime) + int64(e.app.KeepAliveTime)
		s.addEvent(&AppFinishEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp() + int64(e.app.InitTime) + int64(e.app.KeepAliveTime),
			},
			app:       e.app,
			container: s.AppContainerMap[e.app.AppID],
		})
	}

	if e.function != nil {
		e.app.FunctionCnt += 1
		startTime := e.getTimestamp()
		if e.app.InitDoneTimeStamp > e.getTimestamp() {
			startTime = e.app.InitDoneTimeStamp
		}
		//! appInit -> functionStart
		s.addEvent(&FunctionStartEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: startTime,
			},
			function:  e.function,
			app:       e.app,
			container: s.AppContainerMap[e.app.AppID],
		})
	}
}

func (s *Server) handleAppFinishEvent(e *AppFinishEvent) { // 销毁容器
	if e.app == nil || e.app.FunctionCnt != 0 || e.app.FinishTime != e.getTimestamp() {
		return
	}
	s.TimeRunningUsage += e.app.RunningGain
	s.MEMRunningUsage += e.app.MemRunningGain

	AppRunningMemUsage[e.app.AppID] += e.app.MemRunningGain
	AppRunningTimeUsage[e.app.AppID] += e.app.RunningGain
	AppMemUsage[e.app.AppID] += int64(e.app.MEMResources) * (e.getTimestamp() - e.app.InitTimeStamp)
	AppTimeUsage[e.app.AppID] += e.getTimestamp() - e.app.InitTimeStamp

	s.MemUsage += int64(e.app.MEMResources) * (e.getTimestamp() - e.app.InitTimeStamp)
	s.TimeUsage += e.getTimestamp() - e.app.InitTimeStamp
	s.totalMemUsing -= int64(e.app.MEMResources)

	s.AppContainerMap[e.app.AppID] = nil
	if e.app.PreWarmTime > 0 {
		s.addEvent(&AppInitEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp() + int64(e.app.PreWarmTime),
			},
			function: nil,
			app: &Application{
				AppID:             e.app.AppID,
				MEMResources:      e.app.MEMResources,
				FunctionCnt:       0,
				InitTime:          e.app.InitTime,
				InitTimeStamp:     e.getTimestamp() + int64(e.app.PreWarmTime),
				InitDoneTimeStamp: e.getTimestamp() + int64(e.app.PreWarmTime+e.app.InitTime),
				KeepAliveTime:     e.app.KeepAliveTime,
				PreWarmTime:       0,
				FinishTime:        0,
				Left:              int64(-1),
			},
		})
	}
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
	preTime := int64(0)
	//! batchFunctionSubmit -> functionSubmit
	for _, req := range requests {
		if preTime != 0 {
			interval := float64(req.ArrivalTime - preTime)
			// 向上取整
			interval_min := int(math.Ceil(interval / (1000 * 60)))
			appHistogram[req.AppID].sum += 1
			appHistogram[req.AppID].array[interval_min] += 1
		}
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
		app.FunctionCnt += 1
		if app.InitDoneTimeStamp < e.getTimestamp() { // 说明容器已经初始化完成
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
				PreWarmTime:       defaultPreWarmTime,
				FinishTime:        0,
				Left:              int64(-1),
			},
		})
	}
}

func RemoveIdleContainer(cont *Container) {
	for i, v := range ContainerIdleList {
		if v == cont {
			if i == len(ContainerIdleList)-1 { // 考虑如果是最后一个元素的情况
				ContainerIdleList = ContainerIdleList[:i]
			} else {
				ContainerIdleList = append(ContainerIdleList[:i], ContainerIdleList[i+1:]...)
			}
			return
		}
	}
}

func getWindow(app *Application) (int, int) {
	prewarmWindow, keepAliveWindow := 0, 0
	sum1, sum2 := 0, 0
	for i := 0; i < 360; i++ {
		if prewarmWindow != 0 {
			sum1 += appHistogram[app.AppID].array[i]
		}
		sum2 += appHistogram[app.AppID].array[i]
		if float64(sum1) >= 0.05*float64(appHistogram[app.AppID].sum) {
			prewarmWindow = i
			if float64(sum1) >= 0.1*float64(appHistogram[app.AppID].sum) {
				prewarmWindow = 0
			}
		}
		if float64(sum2) >= 0.9*float64(appHistogram[app.AppID].sum) {
			keepAliveWindow = i
			break
		}
	}
	return prewarmWindow, keepAliveWindow
}
