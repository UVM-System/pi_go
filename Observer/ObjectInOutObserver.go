package Observer

import "fmt"

type ObjectInOutObserver struct {
	InCallBack func()
	OutCallBack func()
}

// 光电传感器（红外感应）
var ObjectObserver = ObjectInOutObserver{
	OutCallBack: ObjectOutCallBack,
	InCallBack: func() {
		fmt.Println("...in...")
	},
}

func (ob *ObjectInOutObserver) ObjectIn()  {
	ob.InCallBack()
}

func (ob *ObjectInOutObserver) ObjectOut()  {
	ob.OutCallBack()
}

func (ob ObjectInOutObserver) RecvSerialStr(recvStr string)  {
	//实现Observer接口
	//if Door.Status == Opened {   // 只有在开门状态下，才检测物体进出
		switch recvStr {
		case "3":
			fmt.Println("heard 3 ... ")
			ob.ObjectIn()
			break
		case "4":
			fmt.Println("heard 4 ... ")
			ob.ObjectOut()
			break
		}
	//}
}

