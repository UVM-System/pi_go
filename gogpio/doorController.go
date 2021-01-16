package gogpio

import (
	"pi_go/config"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	doorPin rpio.Pin
	mutex   sync.Mutex
)

// 初始化
func init() {
	if err := rpio.Open(); err != nil {
		panic(err.Error())
	}
	doorPin = rpio.Pin(config.Config.DoorPin)
	mutex.Lock()
	doorPin.High()
	mutex.Unlock()
	doorPin.Output()
}

// 开门，等待 delayTime
func OpenDoor(delayTime int) {
	mutex.Lock()
	doorPin.Low()
	mutex.Unlock()
	time.Sleep(time.Second * time.Duration(delayTime))
	mutex.Lock()
	doorPin.High()
	mutex.Unlock()
}

// 获取门控串口状态
func GetDoorPinState() uint8 {
	mutex.Lock()
	defer mutex.Unlock()
	return uint8(doorPin.Read())
}
