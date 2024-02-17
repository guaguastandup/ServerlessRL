package main

type Function struct {
	FuncID   string
	FuncType FUNCTION_TYPE
	AppID    string
	RunTime  int // ms
}

type Application struct {
	AppID        string
	MEMResources int
	// Function     *Function
	FunctionCnt int
	// * time
	InitTime          int // ms
	InitTimeStamp     int64
	InitDoneTimeStamp int64

	KeepAliveTime int
	LastIdleTime  int64
	FinishTime    int64

	// *
	Left        int64
	RunningGain int64
}
