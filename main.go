package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"pi_go/Observer"
	"pi_go/config"
	"pi_go/gogpio"
	"pi_go/serialHandler"
	"sync"
	"time"
	"github.com/stianeikeland/go-rpio"
	"github.com/tarm/serial"
)

var (
	inputcode string
)

func doorEven() {
	fmt.Printf("door closed!")
}

func postImage() {
	Observer.PostAllImage("start")
}

func main() {
	// 等待 5s 等所有摄像头都打开并开始拍照
	time.Sleep(10 * time.Second)
	// 每次开启时，都需要拍一个冰箱的初始状态，记录初始时拥有多少商品
	//postImage()
	// 开始检测串口信息
	go serialandCapPost()
	for order := 0; ; order = 0 {
		fmt.Println("Please input the order: ")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("|   1. Make all cap take pictures and upload --start  |")
		fmt.Println("|   2. Make all cap take pictures and upload --end    |")
		fmt.Println("|   3. Test serial                                    |")
		fmt.Println("|   4. Open the door                                  |")
		fmt.Println("|   5. test gpio                                      |")
		fmt.Println("|   6. Start                                          |")
		fmt.Println("|   0. Exit                                           |")
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Scanln(&order)
		switch order {
		case 1:
			Observer.PostAllImage("start")
		case 2:
			Observer.PostAllImage("end")
		case 3:
			testserial()
			fmt.Println("Test serial")
		case 4:
			fmt.Println("Open the door")
			gogpio.OpenDoor(3)
			Observer.DoorOpenedCallBack()
		case 5:
			testGPIO()
		case 6:
			start()
		case 0:
			os.Exit(0)
		default:
			fmt.Println("Input error!!! Please input again")
		}
	}
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func start() {
	r := gin.Default()
	r.POST("/openDoor", func(c *gin.Context) {
		data := User{}
		c.BindJSON(&data)
		log.Printf("%v", &data)
		gogpio.OpenDoor(3)
		Observer.DoorOpenedCallBack()
		c.JSON(http.StatusOK, gin.H{
			"username": data.Username,
			"password": data.Password,
		})
	})
	r.Run(":8000") // listen and serve on 0.0.0.0:8080
}

func serialandCapPost() {
	var wg sync.WaitGroup
	serialObservers := []serialHandler.SerialObserver{Observer.ObjectObserver, &Observer.Door}
	wg.Add(1)
	go serialHandler.SerialPortListen(config.Config.SerialPort, config.Config.Baudrate, serialObservers)
	for {
		responseStr := <-Observer.ObjectOutChanel
		log.Print(responseStr)
	}
	wg.Wait()
}

func testserial() {
	c := &serial.Config{
		Name: config.Config.SerialPort,
		Baud: config.Config.Baudrate,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		str := string(buf[:n])
		log.Print(str)
		//log.Printf("%q", buf[:n])
	}
}

func testGPIO() {
	var pin = rpio.Pin(16)
	if err := rpio.Open(); err != nil {
		os.Exit(1)
	}
	pin.High()
	pin.Output()
	order := 0
	for {
		order = 0
		fmt.Println("input order...")
		fmt.Scanln(&order)
		switch order {
		case 1:
			pin.High()
		case 2:
			pin.Low()
		case 0:
			return
		}
	}
}

func getTimestr() string {
	timeStr := time.Now().Format("2006-01-02,15:04:05") //当前时间的字符串，2006-01-02 15:04:05(06年1月2号下午3点4分5秒，06 12345 不重复)，固定写法
	return timeStr
}
