package Observer

import (
	"fmt"
	"log"
	"pi_go/gogpio"
)

type DoorStatus int

const (
	_ DoorStatus = iota
	Opened
	Closed
)
type DoorObserver struct {
	Status DoorStatus
}

// 一个门
var Door = DoorObserver{Status:Closed}

//TODO:
func (doorOb *DoorObserver) Open()  {
	//观察到开门, 设置开门时的初始状态
	DoorOpenedCallBack()
	log.Println("the door is opened...")
}
//TODO:
func (doorOb *DoorObserver) Close() {
	//观察到门已经被关上,应该开始结算
	log.Println("the door is closed...")
}
func (doorOb *DoorObserver) RecvSerialStr(recvStr string)  {
	//实现Observer接口
	switch recvStr {
	case "5":
		// 只有在门控串口信号为低电位时，才可以开门
		if gogpio.GetDoorPinState() == 0 {
			doorOb.Status = Opened
			fmt.Println("heard 5 ... ")
			doorOb.Open()
		}
	case "6":
		fmt.Println("heard 6 ... ")
		doorOb.Status = Closed
		doorOb.Close()
	}
}





