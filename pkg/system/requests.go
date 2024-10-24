package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// var Prefix string = "/Users/zhangxinyue/go/src/serverlessRL/dataset/azurefunctions/"
var Prefix string = "/xinyue/ServerlessRL/dataset/azurefunctions/"

type FUNCTION_TYPE int

const (
	Http FUNCTION_TYPE = iota
	Queue
	Event
	Orchestration
	Timer
	Storage
	Others
)

var FunctionTypeMap = map[string]FUNCTION_TYPE{
	"http":          Http,
	"queue":         Queue,
	"event":         Event,
	"orchestration": Orchestration,
	"timer":         Timer,
	"storage":       Storage,
	"others":        Others,
}

type Request struct { // 对N种job进行测量即可
	ID       int
	AppID    string // 表示运行的任务是什么
	FuncID   string
	FuncType FUNCTION_TYPE
	// * time
	ArrivalTime int64 // 任务的到达时间戳, 单位为秒
	LoadTime    int   // time 1. 任务加载/虚拟机创建的时间, 不包括镜像拉取的时间,
	RunTime     int   // time 2. 任务运行的时间
	// * resource
	MEMResources float64
}

var RequestID int = 0
var MemoryMap map[string]int = make(map[string]int)                         // appID
var MemoryFuncMap map[string]int = make(map[string]int)                     // appID -> func memory
var DurationMap map[string]map[string]int = make(map[string]map[string]int) // appID -> functionID -> duration
var ColdStartTimeMap map[string]int = make(map[string]int)                  // appID -> cold start time
var AppFuncCntMap map[string]int = make(map[string]int)                     // appID -> function cnt
var ColdStartCntMap map[string]int = make(map[string]int)

var NewMemoryMap map[string]int = make(map[string]int) // appID

func initMap() {
	MemoryMap = make(map[string]int)
	MemoryFuncMap = make(map[string]int)
	DurationMap = make(map[string]map[string]int)
	ColdStartTimeMap = make(map[string]int)
	AppFuncCntMap = make(map[string]int)
	ColdStartCntMap = make(map[string]int)
}

func ParseAppFuncCnt(day int) {
	var duration_file string
	if day < 10 {
		duration_file = fmt.Sprintf(Prefix+"duration/duration_d0%d.csv", day)
	} else {
		duration_file = fmt.Sprintf(Prefix+"duration/duration_d%d.csv", day)
	}
	file, err := os.Open(duration_file)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 对于每个app, 记录它有多少个func
	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置字段分隔符为空格，因为你的CSV数据以空格分隔
	// 读取所有的数据
	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(allRecords); i++ {
		appID := allRecords[i][0]
		AppFuncCntMap[appID]++
	}
}

func ParseMemory(day int) {
	var memory_file string
	if day < 10 {
		memory_file = fmt.Sprintf(Prefix+"memory/mem_d0%d.csv", day)
	} else {
		memory_file = fmt.Sprintf(Prefix+"memory/mem_d%d.csv", day)
	}
	file, err := os.Open(memory_file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置字段分隔符为空格，因为你的CSV数据以空格分隔

	// 读取所有的数据
	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	ParseAppFuncCnt(day)

	maxMemory := 0
	for i := 1; i < len(allRecords); i++ {
		appID := allRecords[i][0]
		memory, err := strconv.Atoi(allRecords[i][1])
		if err != nil {
			fmt.Println("error memory: ", memory)
			panic(err)
		}
		MemoryMap[appID] = memory
		if memory > maxMemory {
			maxMemory = memory
		}
	}
	for k, v := range MemoryMap {
		MemoryFuncMap[k] = v / AppFuncCntMap[k]
	}
	fmt.Println("maxMemory: ", maxMemory)
}

func ParseDuration(day int) {
	var duration_file string
	if day < 10 {
		duration_file = fmt.Sprintf(Prefix+"duration/duration_d0%d.csv", day)
	} else {
		duration_file = fmt.Sprintf(Prefix+"duration/duration_d%d.csv", day)
	}
	fmt.Println("duration_file: ", duration_file)
	file, err := os.Open(duration_file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置字段分隔符为空格，因为你的CSV数据以空格分隔

	// 读取所有的数据
	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for i := 1; i < len(allRecords); i++ {
		appID := allRecords[i][0]
		functionID := allRecords[i][1]
		duration, err := strconv.Atoi(allRecords[i][2])
		if duration < 0 {
			continue
		}
		if err != nil {
			panic(err)
		}
		if _, ok := DurationMap[appID]; !ok {
			DurationMap[appID] = make(map[string]int)
		}
		DurationMap[appID][functionID] = duration

		ColdStartTimeMap[appID] += duration
		ColdStartCntMap[appID]++
	}
	for k, v := range ColdStartTimeMap {
		ColdStartTimeMap[k] = v / ColdStartCntMap[k]
	}
}

func ParseRequests(day int, minute int) []*Request {
	// 读取csv文件
	var requests []*Request = make([]*Request, 0)
	var invocation_file string
	if day < 10 {
		// invocation_file = fmt.Sprintf("/Users/zhangxinyue/go/src/serverlessRL/dataset/fake/invocation/d0%d/invocation_d0%d_m%d.csv", day, day, minute)
		invocation_file = fmt.Sprintf(Prefix+"invocation/d0%d/invocation_d0%d_m%d.csv", day, day, minute)
	} else {
		// invocation_file = fmt.Sprintf("/Users/zhangxinyue/go/src/serverlessRL/dataset/fake/invocation/d%d/invocation_d%d_m%d.csv", day, day, minute)
		invocation_file = fmt.Sprintf(Prefix+"invocation/d%d/invocation_d%d_m%d.csv", day, day, minute)
	}
	file, err := os.Open(invocation_file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置字段分隔符为空格，因为你的CSV数据以空格分隔

	// 读取所有的数据
	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println("allRecords: ", len(allRecords)-1)

	// 遍历除了第一列（name列）以外的每一列
	for i := 1; i < len(allRecords); i++ { // 假设所有行的列数相同
		appID := allRecords[i][0]
		functionID := allRecords[i][1]
		functionType := allRecords[i][2]
		arrivalCnt, err := strconv.ParseFloat(allRecords[i][3], 64)
		if err != nil {
			panic(err)
		}
		if MemoryMap[appID] == 0 || DurationMap[appID][functionID] == 0 {
			continue
		}
		if arrivalCnt > float64(ArricalCnt) {
			arrivalCnt = float64(ArricalCnt)
		}
		interval := int((60 * 1000) / arrivalCnt)
		ColdStartTimeMap[functionID] = ColdStartTimeMap[appID]
		MemoryFuncMap[functionID] = MemoryFuncMap[appID]
		MemoryMap[functionID] = MemoryFuncMap[functionID]
		for j := 0; j < int(arrivalCnt); j++ {
			request := Request{
				// ID:    NewRequestID(),
				// AppID: appID,
				AppID:        functionID,
				FuncID:       functionID,
				FuncType:     FunctionTypeMap[functionType],
				ArrivalTime:  int64(1140*day*Minute + 60*1000*minute + interval*j),
				RunTime:      int(float64(DurationMap[appID][functionID])),
				LoadTime:     int(float64(ColdStartTimeMap[appID]) * 1.5),
				MEMResources: float64(MemoryMap[functionID]),
			}
			if request.RunTime < 0 {
				continue
			}
			requests = append(requests, &request)
		}
	}
	// sort by arrival time
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].ArrivalTime < requests[j].ArrivalTime
	})
	return requests
}

func NewRequestID() int {
	RequestID++
	return RequestID
}
