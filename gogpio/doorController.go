package gogpio

import (
	"github.com/stianeikeland/go-rpio"
	"pi_go/config"
	"sync"
	"time"
)

var (
	doorPin rpio.Pin
	pi2StartPin rpio.Pin
	pi2EndPin rpio.Pin
	mutex sync.Mutex
)
// 初始化
func init() {
	if err := rpio.Open(); err != nil {
		panic(err.Error())
	}
	doorPin = rpio.Pin(config.Config.DoorPin)
	pi2StartPin = rpio.Pin(config.Config.Pi2StartPin)
	pi2EndPin = rpio.Pin(config.Config.Pi2EndPin)
	mutex.Lock()
	doorPin.High()
	pi2StartPin.High()
	pi2EndPin.High()
	mutex.Unlock()
	doorPin.Output()
	pi2StartPin.Output()
	pi2EndPin.Output()
}

// 开门，等待 delayTime
func OpenDoor(delayTime int)  {
	mutex.Lock()
	doorPin.Low()
	mutex.Unlock()
	go Pi2postImages("start")
	time.Sleep(time.Second * time.Duration(delayTime))
	mutex.Lock()
	doorPin.High()
	mutex.Unlock()
}

// 让第二个树莓派发送图片
func Pi2postImages(order string)  {
	if order == "start" {
		mutex.Lock()
		pi2StartPin.Low()
		mutex.Unlock()
		// 每次停止 0.1 s
		time.Sleep(time.Millisecond * 2000)
		mutex.Lock()
		pi2StartPin.High()
		mutex.Unlock()
	} else if order == "end" {
		mutex.Lock()
		pi2EndPin.Low()
		mutex.Unlock()
		// 每次停止 0.1 s
		time.Sleep(time.Millisecond * 100)
		mutex.Lock()
		pi2EndPin.High()
		mutex.Unlock()
	}
}

// 获取门控串口状态
func GetDoorPinState() uint8 {
	mutex.Lock()
	defer mutex.Unlock()
	return uint8(doorPin.Read())
}

func GetPi2Pin() (uint8, uint8) {
	mutex.Lock()
	defer mutex.Unlock()
	return uint8(pi2StartPin.Read()), uint8(pi2EndPin.Read())
}
