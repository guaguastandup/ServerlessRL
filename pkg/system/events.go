package main

import (
	log "github.com/sirupsen/logrus"
)

type FunctionSubmitEvent struct {
	baseEvent
	function *Function
}

type FunctionStartEvent struct {
	baseEvent
	function  *Function
	app       *Application
	container *Container
}

type FunctionFinishEvent struct {
	baseEvent
	function  *Function
	app       *Application
	container *Container
}

type AppInitEvent struct {
	baseEvent
	app      *Application
	function *Function
}

type AppFinishEvent struct {
	baseEvent
	app       *Application
	container *Container
}

type BatchFunctionSubmitEvent struct {
	baseEvent
	day    int
	minute int
}

func (e *FunctionSubmitEvent) log() {
	log.WithFields(log.Fields{
		"event":     FunctionSubmitEventType.String(),
		"timestamp": e.getTimestamp(),
		"function":  e.function.FuncID,
		"app":       e.function.AppID,
		"funcType":  e.function.FuncType,
		"runTime":   e.function.RunTime,
	}).Info(FunctionSubmitEventType)
}

func (e *FunctionSubmitEvent) String() string {
	return FunctionSubmitEventType.String()
}

func (e *FunctionStartEvent) log() {
	log.WithFields(log.Fields{
		"event":     FunctionStartEventType.String(),
		"timestamp": e.getTimestamp(),
		"function":  e.function.FuncID,
		"app":       e.function.AppID,
		"funcType":  e.function.FuncType,
		"runTime":   e.function.RunTime,
		"container": e.container.ID,
	}).Info(FunctionStartEventType)
}

func (e *FunctionStartEvent) String() string {
	return FunctionStartEventType.String()
}

func (e *FunctionFinishEvent) log() {
	log.WithFields(log.Fields{
		"event":     FunctionFinishEventType.String(),
		"timestamp": e.getTimestamp(),
		"function":  e.function.FuncID,
		"app":       e.function.AppID,
		"funcType":  e.function.FuncType,
		"runTime":   e.function.RunTime,
		"container": e.container.ID,
	}).Info(FunctionFinishEventType)
}

func (e *FunctionFinishEvent) String() string {
	return FunctionFinishEventType.String()
}

func (e *AppInitEvent) log() {
	log.WithFields(log.Fields{
		"event":     AppInitEventType.String(),
		"timestamp": e.getTimestamp(),
		"app":       e.app.AppID,
		"mem":       e.app.MEMResources,
		"initTime":  e.app.InitTime,
	}).Info(AppInitEventType)
}

func (e *AppInitEvent) String() string {
	return AppInitEventType.String()
}

func (e *AppFinishEvent) log() {
	log.WithFields(log.Fields{
		"event":     AppFinishEventType.String(),
		"timestamp": e.getTimestamp(),
		"app":       e.app.AppID,
		"container": e.container.ID,
	}).Info(AppFinishEventType)
}

func (e *AppFinishEvent) String() string {
	return AppFinishEventType.String()
}

func (e *BatchFunctionSubmitEvent) log() {
	log.WithFields(log.Fields{
		"event":     BatchFunctionSubmitEventType.String(),
		"timestamp": e.getTimestamp(),
		"day":       e.day,
		"minute":    e.minute,
	}).Info(BatchFunctionSubmitEventType)
}

func (e *BatchFunctionSubmitEvent) String() string {
	return BatchFunctionSubmitEventType.String()
}

type eventType int

const (
	FunctionSubmitEventType eventType = iota
	FunctionStartEventType
	FunctionFinishEventType
	AppInitEventType
	AppFinishEventType
	BatchFunctionSubmitEventType
)

func (e eventType) String() string {
	names := [...]string{
		"FunctionSubmitEvent",
		"FunctionStartEvent",
		"FunctionFinishEvent",
		"AppInitEvent",
		"AppFinishEvent",
		"BatchFunctionSubmitEvent",
	}
	if e < FunctionSubmitEventType || e > BatchFunctionSubmitEventType {
		log.Panic("Unknown event type: " + string(e))
	}
	return names[e]
}
