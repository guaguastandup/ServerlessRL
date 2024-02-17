package main

import "fmt"

func (s *Server) handleFuncSubmitEvent(e *FunctionSubmitEvent) {
	appID := e.function.AppID
	appMemory := MemoryMap[appID]
	s.appCnt[appID] += 1
	s.totalRequest += 1

	if len(s.AppContainerIdleMap[appID]) > 0 { // warm start
		warmStartCnt++
		s.appWarmUpCnt[appID] += 1
		container := s.AppContainerIdleMap[appID][0]
		s.AppContainerIdleMap[appID] = s.AppContainerIdleMap[appID][1:]

		//! functionSubmit -> functionStart
		s.addEvent(&FunctionStartEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			},
			app:       container.App,
			container: container,
			function:  e.function,
		})
	} else { // cold start
		//! functionSubmit -> appInit
		s.addEvent(&AppInitEvent{
			baseEvent: baseEvent{
				id:        s.newEventId(),
				timestamp: e.getTimestamp(),
			},
			app: &Application{
				AppID:         appID,
				MEMResources:  appMemory,
				InitTime:      ColdStartTimeMap[appID],
				Function:      e.function,
				KeepAliveTime: defaultKeepAliveTime,
				InitTimeStamp: e.getTimestamp(),
				FinishTime:    0,
			},
		})
	}
}

func (s *Server) handleFuncStartEvent(e *FunctionStartEvent) {
	e.container.Status = 1
	e.app.Function = e.function

	s.totalMemRunning += e.app.MEMResources

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
	e.container.Status = 0
	e.app.Function = nil
	e.app.LastIdleTime = e.getTimestamp()

	s.MEMRunningUsage += (e.app.MEMResources * e.function.RunTime)
	s.TimeRunningUsage += e.function.RunTime
	s.totalMemRunning -= e.app.MEMResources
	s.AppContainerIdleMap[e.app.AppID] = append(s.AppContainerIdleMap[e.app.AppID], e.container)

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
	s.totalMemUsing += int64(e.app.MEMResources)

	//! appInit -> functionStart
	s.addEvent(&FunctionStartEvent{
		baseEvent: baseEvent{
			id:        s.newEventId(),
			timestamp: e.getTimestamp() + int64(e.app.InitTime),
		},
		function:  e.app.Function,
		app:       e.app,
		container: cont,
	})
}

func (s *Server) handleAppFinishEvent(e *AppFinishEvent) { // 销毁容器
	if e.app == nil {
		return
	}
	if e.app.Function != nil {
		return
	}
	if e.app.FinishTime != e.getTimestamp() {
		return
	}
	s.MemUsage += (e.app.MEMResources * int((e.getTimestamp() - e.app.InitTimeStamp)))
	s.TimeUsage += int(e.getTimestamp() - e.app.InitTimeStamp)
	s.totalMemUsing -= int64(e.app.MEMResources)
	delete(s.ContainerMap, e.container.ID)
	if len(s.AppContainerIdleMap[e.app.AppID]) > 0 {
		for i, v := range s.AppContainerIdleMap[e.app.AppID] {
			if v.ID == e.container.ID {
				if i == len(s.AppContainerIdleMap[e.app.AppID])-1 {
					s.AppContainerIdleMap[e.app.AppID] = s.AppContainerIdleMap[e.app.AppID][:i]
				} else {
					s.AppContainerIdleMap[e.app.AppID] = append(s.AppContainerIdleMap[e.app.AppID][:i], s.AppContainerIdleMap[e.app.AppID][i+1:]...)
				}
				break
			}
		}
	}
	e.app = nil
}

func (s *Server) handleBatchFuncSubmitEvent(e *BatchFunctionSubmitEvent) {
	fmt.Printf("Batch Submit Event: %d day %d minute\n\n", e.day, e.minute)
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
