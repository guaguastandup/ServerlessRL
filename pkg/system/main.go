package main

import (
	"container/heap"
	"fmt"
	"log"
	"time"
)

var Second int = 1000 // ms
var Minute int = 60 * Second
var defaultKeepAliveTime int = 5 * Minute

type Container struct {
	ID               int
	App              *Application
	MEMUsage         int // 申请的Usage
	MEMRunningUsage  int // 运行时的Usage
	TimeUsage        int // 申请的Usage
	TimeRunningUsage int // 运行时的Usage
}

type Server struct { // Server-wide
	// * event
	EventQueue  eventQueue
	currEventId int
	// * Container Map
	AppContainerMap  map[string]*Container
	totalContainerID int // 用于生成ContainerID
	// * usage
	MEMCapacity      int
	MemUsage         int // 当前的MEM总使用量
	MEMRunningUsage  int // 任务运行的时候的MEM使用率, 有效值, unitMEMUsage * unitTime
	TimeUsage        int
	TimeRunningUsage int
	// * Time
	currTime        int64
	totalMemUsing   int64
	totalMemRunning int
	// * statistics
	warmStartCnt    int
	totalRequest    int64
	appWarmStartCnt map[string]int
	appRequestCnt   map[string]int
}

func (s *Server) Run() {
	start := time.Now()
	for s.EventQueue.Len() > 0 {
		e := heap.Pop(&s.EventQueue).(event)
		if e.getTimestamp() < s.currTime {
			e.log()
			fmt.Println(e.getTimestamp(), s.currTime)
			panic("Event is not in chronological order")
		}
		s.currTime = e.getTimestamp()
		s.handleEvent(e)
		if e.String() == "BatchFunctionSubmitEvent" {
			fmt.Printf("MemOccupyingUsage: %.1f GB\n", float64(s.totalMemUsing/1024.0))
			fmt.Printf("MEMRunningUsage: %.1f GB\n", float64(s.totalMemRunning/1024.0))
			fmt.Printf("Mem Score: %.4f %%\n", 100.0*float64(s.MEMRunningUsage)/float64(s.MemUsage))
			fmt.Printf("Time Score: %.4f %%\n", 100.0*float64(s.TimeRunningUsage)/float64(s.TimeUsage))
			fmt.Printf("warmStart Rate: %.4f %%\n\n", 100.0*float64(s.warmStartCnt)/float64(s.totalRequest))
		}
	}
	fmt.Printf("Simulation takes %v", time.Since(start))
}

func (s *Server) handleEvent(e event) {
	switch v := e.(type) {
	case *FunctionSubmitEvent:
		s.handleFuncSubmitEvent(v)
	case *FunctionStartEvent:
		s.handleFuncStartEvent(v)
	case *FunctionFinishEvent:
		s.handleFuncFinishEvent(v)
	case *AppInitEvent:
		s.handleAppInitEvent(v)
	case *AppFinishEvent:
		s.handleAppFinishEvent(v)
	case *BatchFunctionSubmitEvent:
		s.handleBatchFuncSubmitEvent(v)
	default:
		log.Panic("Unknown event type")
	}
}

func main() {
	Server := &Server{
		MEMCapacity:     1024 * 10,
		AppContainerMap: make(map[string]*Container),
		appWarmStartCnt: make(map[string]int),
		appRequestCnt:   make(map[string]int),
	}
	for day := 1; day <= 1; day++ {
		for i := 0; i < 1140; i++ {
			Server.addEvent(&BatchFunctionSubmitEvent{
				baseEvent: baseEvent{
					id:        Server.newEventId(),
					timestamp: int64(day*1140*Minute + i*Minute),
				},
				day:    day,
				minute: i + 1,
			})
		}
	}
	Server.Run()
	for k, v := range Server.appWarmStartCnt {
		fmt.Println("warmstart rate: ", float64(v)/float64(Server.appRequestCnt[k]))
	}
}

// Constructor
func (s *Server) NewContainer(e *AppInitEvent) (cont *Container) {
	cont = &Container{
		ID:  s.totalContainerID,
		App: e.app,
	}
	s.totalContainerID++
	s.AppContainerMap[e.app.AppID] = cont
	return cont
}
