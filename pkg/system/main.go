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

var defaultKeepAliveTime int = 5 * Minute
var defaultPreWarmTime int = 1 * Minute
var defaultMemoryCapcity int = 1024 * 1024 // 1TB
var ArricalCnt int = 5
var IsFixed int = 0
var SumLimit int = 50
var leftBound float64 = 0.05
var leftBound2 float64 = 0.1
var rightBound float64 = 0.95

var policy string = "lru"

var totalDay int = 1
var totalMinute int = 500

type Container struct {
	ID               int64
	App              *Application
	MEMUsage         int64 // 申请的Usage
	MEMRunningUsage  int64 // 运行时的Usage
	TimeUsage        int64 // 申请的Usage
	TimeRunningUsage int64 // 运行时的Usage
	Index            int
}

type Server struct { // Server-wide
	// * event
	EventQueue  eventQueue
	currEventId int64
	// * Container Map
	AppContainerMap  map[string]*Container
	totalContainerID int64 // 用于生成ContainerID
	// * usage
	MEMCapacity      int64
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
	heap.Init(h) // 初始化堆
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
			fmt.Printf("Evicted Memory: %.1f GB\n", float64(EvictedMemory/1024.0))
			fmt.Printf("warmStart Rate: %.4f %%\n\n", 100.0*float64(s.warmStartCnt)/float64(s.totalRequest))
			if e.(*BatchFunctionSubmitEvent).minute%20 == 0 {
				sum := float64(0)
				cnt := 0
				for k, v := range s.appRequestCnt {
					sum += (1.0 - float64(s.appWarmStartCnt[k])/float64(v))
					cnt += 1
				}
				fmt.Printf("app average coldstart rate: %.4f %%\n", float64(sum)/float64(cnt))
			}
			startTime = time.Now()
			cnt = 0
			if e.(*BatchFunctionSubmitEvent).minute%100 == 0 {
				sum_mem_score, sum_time_score := 0.0, 0.0
				cnt_mem_score, cnt_time_score := 0, 0
				for k, v := range AppMemUsage {
					sum_mem_score += 100.0 * float64(AppRunningMemUsage[k]) / float64(v)
					cnt_mem_score += 1
				}
				for k, v := range AppTimeUsage {
					sum_time_score += 100.0 * float64(AppRunningTimeUsage[k]) / float64(v)
					cnt_time_score += 1
				}
				fmt.Printf("average mem socre: %.5f\n", sum_mem_score/float64(cnt_mem_score))
				fmt.Printf("average time socre: %.5f\n", sum_time_score/float64(cnt_time_score))
			}
			if e.(*BatchFunctionSubmitEvent).minute == totalMinute && e.(*BatchFunctionSubmitEvent).day == totalDay {
				sum_mem_score, sum_time_score := 0.0, 0.0
				cnt_mem_score, cnt_time_score := 0, 0
				for k, v := range AppMemUsage {
					fmt.Printf("app mem socre: %.5f\n", 100.0*float64(AppRunningMemUsage[k])/float64(v))
					sum_mem_score += 100.0 * float64(AppRunningMemUsage[k]) / float64(v)
					cnt_mem_score += 1
				}
				for k, v := range AppTimeUsage {
					fmt.Printf("app time socre: %.5f\n", 100.0*float64(AppRunningTimeUsage[k])/float64(v))
					sum_time_score += 100.0 * float64(AppRunningTimeUsage[k]) / float64(v)
					cnt_time_score += 1
				}
				fmt.Printf("average mem socre: %.5f\n", sum_mem_score/float64(cnt_mem_score))
				fmt.Printf("average time socre: %.5f\n", sum_time_score/float64(cnt_time_score))
			}
		}
		s.handleEvent(e)
	}

	fmt.Printf("MemOccupyingUsage: %.1f GB\n", float64(s.totalMemUsing/1024.0))
	fmt.Printf("MEMRunningUsage: %.1f GB\n", float64(s.totalMemRunning/1024.0))
	fmt.Printf("Mem Score: %.4f %%\n", 100.0*float64(s.MEMRunningUsage)/float64(s.MemUsage))
	fmt.Printf("Time Score: %.4f %%\n", 100.0*float64(s.TimeRunningUsage)/float64(s.TimeUsage))
	fmt.Printf("warmStart Rate: %.4f %%\n\n", 100.0*float64(s.warmStartCnt)/float64(s.totalRequest))
	fmt.Printf("totalRequest: %d\n", s.totalRequest)
	fmt.Printf("warmStartCnt: %d\n", s.warmStartCnt)
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
	num, _ := strconv.ParseFloat(os.Args[1], 64)
	defaultKeepAliveTime = int(num * float64(Minute))
	num, _ = strconv.ParseFloat(os.Args[2], 64)
	defaultPreWarmTime = int(num * float64(Minute))
	num, _ = strconv.ParseFloat(os.Args[3], 64)
	defaultMemoryCapcity = int(num * 1024) // num GB
	num, _ = strconv.ParseFloat(os.Args[4], 64)
	ArricalCnt = int(num)
	num, _ = strconv.ParseFloat(os.Args[5], 64)
	IsFixed = int(num)
	num, _ = strconv.ParseFloat(os.Args[6], 64)
	SumLimit = int(num)
	num, _ = strconv.ParseFloat(os.Args[7], 64)
	leftBound = num
	num, _ = strconv.ParseFloat(os.Args[8], 64)
	leftBound2 = num
	num, _ = strconv.ParseFloat(os.Args[9], 64)
	rightBound = num

	policy = os.Args[10]

	fmt.Printf("default KeepAliveTime: %d\n", defaultKeepAliveTime)
	fmt.Printf("default PreWarmTime: %d\n", defaultPreWarmTime)
	fmt.Printf("default MemoryCapcity: %d\n", defaultMemoryCapcity)
	fmt.Printf("ArricalCnt: %d\n", ArricalCnt)
	fmt.Printf("IsFixed: %d\n", IsFixed)
	fmt.Printf("SumLimit: %d\n", SumLimit)
	fmt.Printf("leftBound: %.2f\n", leftBound)
	fmt.Printf("leftBound2: %.2f\n", leftBound2)
	fmt.Printf("rightBound: %.2f\n", rightBound)
	fmt.Printf("policy: %v\n", policy)

	Server := &Server{
		MEMCapacity:     int64(defaultMemoryCapcity),
		AppContainerMap: make(map[string]*Container),
		appWarmStartCnt: make(map[string]int),
		appRequestCnt:   make(map[string]int),
	}

	for day := 1; day <= totalDay; day++ {
		for i := 0; i < totalMinute; i++ {
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
