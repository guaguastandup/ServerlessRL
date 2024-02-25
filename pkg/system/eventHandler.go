package main

var AppFuncMap map[string]map[string]int = make(map[string]map[string]int)  // appID -> functionID -> start time
var AppLeft map[string]map[string]int64 = make(map[string]map[string]int64) // appID -> left time
var AppMemUsage map[string]int64 = make(map[string]int64)                   // appID -> mem usage
var AppRunningMemUsage map[string]int64 = make(map[string]int64)            // appID -> running mem usage
var AppTimeUsage map[string]int64 = make(map[string]int64)                  // appID -> time usage
var AppRunningTimeUsage map[string]int64 = make(map[string]int64)           // appID -> running time usage

var LastIdleTime map[string]int64 = make(map[string]int64) // appID -> last idle time

var KeepAliveTimeMap map[string]int = make(map[string]int) // appID -> keep alive time

func (s *Server) handleFuncStartEvent(e *FunctionStartEvent) {
	if IsExistInIdleList(e.container) {
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
	if e.app.FinishTime <= e.getTimestamp()+int64(e.function.RunTime) { // 为了避免函数没有执行结束, App就被销毁的情况
		e.app.FinishTime = e.getTimestamp() + int64(e.function.RunTime) + 1
	}
	//! Function Start -> Function Finish
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

	if e.container == nil {
		panic("container is nil")
	}
	if e.app.FunctionCnt == 0 && !IsExistInIdleList(e.container) {
		s.AddToIdleList(e.container)
		LastIdleTime[e.app.AppID] = e.getTimestamp()
	}

	prewarmWindow, keepAliveWindow := getWindow(e.app)
	e.app.KeepAliveTime = keepAliveWindow
	e.app.PreWarmTime = prewarmWindow
	KeepAliveTimeMap[e.app.AppID] = keepAliveWindow

	if prewarmWindow == 0 { // 使用KeepAlive的策略
		//! functionFinish -> appFinish
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
		if e.app.FinishTime <= e.getTimestamp() {
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
}

func (s *Server) handleAppInitEvent(e *AppInitEvent) { // 冷启动
	flag := 0
	if s.AppContainerMap[e.app.AppID] == nil {
		s.totalMemUsing += int64(e.app.MEMResources)
		if s.totalMemUsing > s.MEMCapacity { // Memory Overload
			s.handleEvictEvent(&baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			})
		}
		cont := s.NewContainer(e)
		s.AppContainerMap[e.app.AppID] = cont
		flag = 1
		e.app.InitTimeStamp = e.getTimestamp()
		e.app.InitDoneTimeStamp = e.getTimestamp() + int64(e.app.InitTime)
		AppFuncMap[e.app.AppID] = make(map[string]int)
		AppLeft[e.app.AppID] = make(map[string]int64)
	}
	if e.app == nil {
		panic("impossible: nil app")
	}
	if s.AppContainerMap[e.app.AppID] == nil {
		panic("impossible: nil container")
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
		e.app.FunctionCnt += 1 // 表示这个app已经被预定了
		//! appInit -> functionStart
		startTime := e.getTimestamp()
		if e.app.InitDoneTimeStamp > e.getTimestamp() {
			startTime = e.app.InitDoneTimeStamp + 1
		}
		if startTime == e.getTimestamp() {
			s.handleFuncStartEvent(&FunctionStartEvent{
				baseEvent: baseEvent{
					id:        s.newEventId(),
					timestamp: startTime,
				},
				function:  e.function,
				app:       e.app,
				container: s.AppContainerMap[e.app.AppID],
			})
		} else {
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
	} else {
		if !IsExistInIdleList(s.AppContainerMap[e.app.AppID]) { // 如果当前只是预热, 没有明确的Function, 那么就把App加入到IdleList
			s.AddToIdleList(s.AppContainerMap[e.app.AppID])
			LastIdleTime[e.app.AppID] = e.getTimestamp()
		}
	}
}

func (s *Server) handleAppFinishEvent(e *AppFinishEvent) { // 销毁容器
	if e.app == nil || e.container == nil || s.AppContainerMap[e.app.AppID] == nil || e.app.FunctionCnt != 0 || e.app.FinishTime != e.getTimestamp() {
		return
	}
	if IsExistInIdleList(e.container) {
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

	delete(s.AppContainerMap, e.app.AppID)
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
				Score:             s.getScore(e.app.AppID, e.getTimestamp()),
			},
		})
	}
	e.app = nil
	e.container = nil
}

func (s *Server) handleFuncSubmitEvent(e *FunctionSubmitEvent) {
	appID := e.function.AppID
	appMemory := MemoryMap[appID]

	s.appRequestCnt[appID] += 1
	s.totalRequest += 1

	if s.AppContainerMap[appID] != nil { // warm start
		container := s.AppContainerMap[appID]
		startTime := e.getTimestamp()
		if container.App.InitDoneTimeStamp <= e.getTimestamp() { // 说明容器已经初始化完成
			s.warmStartCnt++
			s.appWarmStartCnt[appID] += 1
		} else {
			startTime = container.App.InitDoneTimeStamp + 1
		}
		container.App.FunctionCnt += 1
		//! functionSubmit -> functionStart
		if startTime == e.getTimestamp() {
			s.handleFuncStartEvent(&FunctionStartEvent{
				baseEvent: baseEvent{
					id:        s.newEventId(),
					timestamp: startTime,
				},
				app:       container.App,
				container: container,
				function:  e.function,
			})
		} else {
			s.addEvent(&FunctionStartEvent{
				baseEvent: baseEvent{
					id:        s.newEventId(),
					timestamp: startTime,
				},
				app:       container.App,
				container: container,
				function:  e.function,
			})
		}
	} else { // cold start
		//! functionSubmit -> appInit
		s.handleAppInitEvent(&AppInitEvent{ // todo: 此处需要考虑delayed hits的问题
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
				Score:             s.getScore(appID, e.getTimestamp()),
			},
		})
	}
}
