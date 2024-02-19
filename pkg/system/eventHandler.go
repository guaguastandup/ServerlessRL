package main

import (
	"fmt"
	"math"
	"math/rand"
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

var preTime map[string]int64 = make(map[string]int64)

type histogram struct {
	sum   int
	array []int
}

var appHistogram map[string]*histogram = make(map[string]*histogram)

func stringEquals(a, b string) bool {
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func (s *Server) handleEvictEvent(e *baseEvent) {
	for s.totalMemUsing > s.MEMCapacity {
		if len(ContainerIdleList) == 0 {
			panic("No idle container to evict!")
		}
		var container *Container
		// 删除第一个空闲的容器
		if stringEquals(policy, "lru") {
			container = ContainerIdleList[0]
		} else if stringEquals(policy, "random") {
			index := rand.Intn(len(ContainerIdleList))
			container = ContainerIdleList[index]
		} else if stringEquals(policy, "maxmem") {
			maxMem := 0
			for _, cont := range ContainerIdleList {
				if cont.App.MEMResources > maxMem {
					maxMem = cont.App.MEMResources
					container = cont
				}
			}
		} else if stringEquals(policy, "maxKeepAlive") {
			maxKeepAlive := 0
			for _, cont := range ContainerIdleList {
				if cont.App.KeepAliveTime > maxKeepAlive {
					maxKeepAlive = cont.App.KeepAliveTime
					container = cont
				}
			}
		} else if stringEquals(policy, "minUsage") {
			minUsage := int64(1e18)
			for _, cont := range ContainerIdleList {
				if cont.App.MemRunningGain < minUsage {
					minUsage = cont.App.MemRunningGain
					container = cont
				}
			}
		} else if stringEquals(policy, "maxColdStartRate") {
			maxColdStart := float64(0.0)
			for _, cont := range ContainerIdleList {
				if float64(s.appWarmStartCnt[cont.App.AppID])/float64(s.appRequestCnt[cont.App.AppID]) >= maxColdStart {
					maxColdStart = float64(s.appWarmStartCnt[cont.App.AppID]) / float64(s.appRequestCnt[cont.App.AppID])
					container = cont
				}
			}
		} else {
			fmt.Println("policy: ", policy)
			panic("Invalid policy! " + policy)
		}
		if !ContainerIdleMap[container] {
			panic("impossible " + policy)
		}
		ContainerIdleMap[container] = false
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
	prewarmWindow, keepAliveWindow := getWindow(e.app)
	e.app.KeepAliveTime = keepAliveWindow
	e.app.PreWarmTime = prewarmWindow
	if prewarmWindow == 0 { // 使用KeepAlive的策略
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
	} else {
		e.app.FinishTime = e.getTimestamp() // 立刻删除
		//! functionFinish -> appFinish
		s.handleAppFinishEvent(&AppFinishEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
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

	if ContainerIdleMap[e.container] {
		ContainerIdleMap[e.container] = false
		RemoveIdleContainer(e.container)
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
	//! batchFunctionSubmit -> functionSubmit
	for _, req := range requests {
		if appHistogram[req.AppID] == nil {
			appHistogram[req.AppID] = &histogram{
				sum:   0,
				array: make([]int, 130),
			}
		}
		if preTime[req.AppID] != 0 {
			interval := float64(req.ArrivalTime - preTime[req.AppID])
			interval_min := int(math.Ceil(interval / (1000 * 60)))
			if interval_min > 120 {
				interval_min = 119
			}
			appHistogram[req.AppID].sum += 1
			appHistogram[req.AppID].array[interval_min] += 1
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
	if IsFixed > 0 {
		return defaultPreWarmTime, defaultKeepAliveTime
	}
	if appHistogram[app.AppID].sum <= SumLimit {
		return defaultPreWarmTime, defaultKeepAliveTime
	}
	prewarmWindow, keepAliveWindow := 0, 0
	sum1, sum2 := 0, 0
	for i := 0; i < 120; i++ {
		if prewarmWindow != 0 {
			sum1 += appHistogram[app.AppID].array[i]
		}
		sum2 += appHistogram[app.AppID].array[i]
		if float64(sum1) >= leftBound*float64(appHistogram[app.AppID].sum) {
			prewarmWindow = i
			if float64(sum1) >= leftBound2*float64(appHistogram[app.AppID].sum) {
				prewarmWindow = 0
			}
		}
		if float64(sum2) >= rightBound*float64(appHistogram[app.AppID].sum) {
			keepAliveWindow = i
			break
		}
	}
	return prewarmWindow * 60 * 1000, keepAliveWindow * 60 * 1000
}

// 0110101010100001000010101
// 0010101010001111001011111
// 0101111111011101011101011
// 1011011110101011110101111
