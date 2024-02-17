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
	Function     *Function
	// * time
	InitTime      int // ms
	KeepAliveTime int
	InitTimeStamp int64
	LastIdleTime  int64
	FinishTime    int64
}
