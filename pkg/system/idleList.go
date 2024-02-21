package main

import (
	"container/list"
	"fmt"
)

var ContainerIdleList *list.List = list.New()
var ContainerIdleMap map[int64]*list.Element = make(map[int64]*list.Element) // delete by container id

func RemoveIdleContainer(container *Container) {
	id := container.ID
	nodeToDelete, exists := ContainerIdleMap[id]
	if exists {
		ContainerIdleList.Remove(nodeToDelete) // 从链表中删除该节点
		delete(ContainerIdleMap, id)           // 从映射中删除该节点的指针
		h.RemoveByID(id)
	} else {
		fmt.Println(len(ContainerIdleMap), ContainerIdleList.Len())
		fmt.Println(container.ID)
		panic("delete no element, " + fmt.Sprintf("%d", container.ID))
	}
}

func (s *Server) AddToIdleList(container *Container) {
	if ContainerIdleMap[container.ID] != nil {
		panic("add an exist element")
	}
	ele := ContainerIdleList.PushBack(container)
	ContainerIdleMap[container.ID] = ele
	container.App.Score = s.getScore(container.App.AppID, s.currTime)
	h.Push(container)
}

func IsExistInIdleList(container *Container) bool {
	return ContainerIdleMap[container.ID] != nil
}

func FrontElement() *Container {
	ele := ContainerIdleList.Front()
	container, ok := ele.Value.(*Container)
	if !ok {
		panic("not a Container element")
	}
	return container
}

func BackElement() *Container {
	ele := ContainerIdleList.Front()
	container, ok := ele.Value.(*Container)
	if !ok {
		panic("not a Container element")
	}
	return container
}

func GetElementByIndex(index int) *Container {
	ele := ContainerIdleList.Front()
	for i := 1; i <= index; i++ {
		ele = ele.Next()
		if i == index {
			container, ok := ele.Value.(*Container)
			if !ok {
				panic("not a Container element")
			}
			return container
		}
	}
	panic("no container satisfied")
}
