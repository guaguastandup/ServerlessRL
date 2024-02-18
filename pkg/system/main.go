package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var Second int = 1000 // ms
var Minute int = 60 * Second

// var defaultKeepAliveTime int = 1 * Second
var defaultKeepAliveTime int = 1
var defaultPreWarmTime int = 1 * Minute

type Container struct {
	ID               int
	App              *Application
	MEMUsage         int64 // 申请的Usage
	MEMRunningUsage  int64 // 运行时的Usage
	TimeUsage        int64 // 申请的Usage
	TimeRunningUsage int64 // 运行时的Usage
}

type Server struct { // Server-wide
	// * event
	EventQueue  eventQueue
	currEventId int64
	// * Container Map
	AppContainerMap  map[string]*Container
	totalContainerID int // 用于生成ContainerID
	// * usage
	MEMCapacity      int
	MemUsage         int64 // 当前的MEM总使用量
	MEMRunningUsage  int64 // 任务运行的时候的MEM使用率, 有效值, unitMEMUsage * unitTime
	TimeUsage        int64
	TimeRunningUsage int64
	// * Time
	currTime        int64
	totalMemUsing   int64
	totalMemRunning int64
	// * statistics
	warmStartCnt    int
	totalRequest    int64
	appWarmStartCnt map[string]int
	appRequestCnt   map[string]int
}

func (s *Server) Run() {
	start := time.Now()
	startTime := time.Now()
	cnt := 0
	for s.EventQueue.Len() > 0 {
		cnt += 1
		e := heap.Pop(&s.EventQueue).(event)
		if e.getTimestamp() < s.currTime {
			e.log()
			fmt.Println(e.getTimestamp(), s.currTime)
			panic("Event is not in chronological order")
		}
		s.currTime = e.getTimestamp()
		if e.String() == "BatchFunctionSubmitEvent" {
			fmt.Println("time cost: ", time.Since(startTime).Seconds())
			fmt.Println("Event Count: ", cnt)
			fmt.Printf("MemOccupyingUsage: %.1f GB\n", float64(s.totalMemUsing/1024.0))
			fmt.Printf("MEMRunningUsage: %.1f GB\n", float64(s.totalMemRunning/1024.0))
			fmt.Printf("Mem Score: %.4f %%\n", 100.0*float64(s.MEMRunningUsage)/float64(s.MemUsage))
			fmt.Printf("Time Score: %.4f %%\n", 100.0*float64(s.TimeRunningUsage)/float64(s.TimeUsage))
			fmt.Printf("warmStart Rate: %.4f %%\n\n", 100.0*float64(s.warmStartCnt)/float64(s.totalRequest))
			startTime = time.Now()
			cnt = 0
		}
		s.handleEvent(e)
	}
	// fmt.Printf("MemOccupyingUsage: %.1f GB\n", float64(s.totalMemUsing/1024.0))
	// fmt.Printf("MEMRunningUsage: %.1f GB\n", float64(s.totalMemRunning/1024.0))
	// fmt.Printf("Mem Score: %.4f %%\n", 100.0*float64(s.MEMRunningUsage)/float64(s.MemUsage))
	// fmt.Printf("Time Score: %.4f %%\n", 100.0*float64(s.TimeRunningUsage)/float64(s.TimeUsage))
	// fmt.Printf("warmStart Rate: %.4f %%\n\n", 100.0*float64(s.warmStartCnt)/float64(s.totalRequest))
	// fmt.Printf("totalRequest: %d\n", s.totalRequest)
	// fmt.Printf("warmStartCnt: %d\n", s.warmStartCnt)
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
	if len(os.Args) > 1 {
		// convert string to float
		num, _ := strconv.ParseFloat(os.Args[1], 64)
		defaultKeepAliveTime = int(num * float64(Minute))
	}
	if len(os.Args) > 2 {
		num, _ := strconv.ParseFloat(os.Args[2], 64)
		defaultPreWarmTime = int(num * float64(Minute))
	}
	fmt.Println("default KeepAliveTime: ", defaultKeepAliveTime)
	fmt.Println("dufault PreWarmTime: ", defaultPreWarmTime)
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
	for k, v := range Server.appRequestCnt {
		fmt.Println("warmstart rate: ", float64(Server.appWarmStartCnt[k])/float64(v))
	}
	for k, v := range AppMemUsage {
		fmt.Printf("app mem socre: %.5f\n", 100.0*float64(AppRunningMemUsage[k])/float64(v))
	}
	for k, v := range AppTimeUsage {
		fmt.Printf("app time socre: %.5f\n", 100.0*float64(AppRunningTimeUsage[k])/float64(v))
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
