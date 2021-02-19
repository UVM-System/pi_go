package camutils

import (
	"fmt"
	"log"
	"pi_go/config"
	"sync"

	"gocv.io/x/gocv"
)

var (
	VideoHandlers []VideoCap
	once          sync.Once
)

func init() {
	fmt.Println("init camera")
	once.Do(InitAndStartCap)
}
func InitAndStartCap() {
	VideoHandlers = make([]VideoCap, 0)
	for i := 0; i < len(config.Config.CapConfigs); i++ {
		videoHandler := VideoCap{
			videoId: config.Config.CapConfigs[i].VideoId,
			img:     gocv.NewMat(),
			mutex:   sync.Mutex{},
			Prefix:  config.Config.CapConfigs[i].Prefix,
		}
		go videoHandler.StartCap()
		// ToDo 加入VideoHandler 前， 存储的 img 不能为空图像
		fmt.Println("add \t", i, "\t cap")
		VideoHandlers = append(VideoHandlers, videoHandler)
	}

}

type VideoCap struct {
	videoId int
	img     gocv.Mat
	mutex   sync.Mutex
	Prefix  string
}

func (cap *VideoCap) GetJpegImageBytes() (buf []byte, err error) {
	cap.mutex.Lock()

	if cap.img.Empty() {
		fmt.Println(cap.videoId)
		fmt.Println("sorry the img of cap is not contained image")
	}
	gocv.IMWrite("404.jpg", cap.img)

	imageBytes, err := gocv.IMEncode(".jpg", cap.img)
	if err != nil {
		log.Print("cap  " + string(cap.videoId) + " error")
		log.Fatal(err.Error())
	}
	cap.mutex.Unlock()
	return imageBytes, err
}

func (cap *VideoCap) StartCap() {
	camHandler, err := gocv.OpenVideoCapture(cap.videoId)
	if err != nil {
		panic(err.Error())
	}
	defer camHandler.Close()
	fmt.Println("videoId:\t", cap.videoId)
	camHandler.Set(gocv.VideoCaptureFrameHeight, 1080)
	camHandler.Set(gocv.VideoCaptureFrameWidth, 1920)
	// fmt.Println("videoId:\t", cap.videoId, "\t", camHandler.Get(gocv.VideoCaptureFrameHeight))
	// fmt.Println("videoId:\t", cap.videoId, "\t", camHandler.Get(gocv.VideoCaptureFrameWidth))
	for {
		cap.mutex.Lock()
		fmt.Println("read picture:\t", cap.videoId)
		camHandler.Read(&cap.img)
		cap.mutex.Unlock()
	}
}

/*
func StartCap(img *gocv.Mat,video_id int,mutex sync.Mutex)  {
	cam_handler,err:=gocv.OpenVideoCapture(video_id)
	if err!=nil{
		panic(err.Error())
	}
	for{
		mutex.Lock()
		cam_handler.Read(img)
		mutex.Unlock()
	}
	defer cam_handler.Close()
}
*/
