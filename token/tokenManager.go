package token

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"pi_go/config"
	"sync"
	"time"
)
// 定时更新自身向服务器更新自身token信息,顺便汇报自身IP信息
// 服务器向本机发送开门请求时候,验证token信息后才能正式开柜门
var (
	token string
	tokenChan chan string
	mutex sync.Mutex
	once sync.Once
)

const (
	//一小时更新一次
	updateTime = 3600
)


type TokenJson struct {
	Message string `json:"message"`
	Token string `json:"token"`
}


func init()  {
	GetTokenHttp()
	once.Do(trueInit)
}
func trueInit()  {
	tokenChan = make(chan string)
	go updateToken()
}

func updateToken()  {
	//定时更新token
	ticker := time.NewTicker(time.Second * updateTime)
	for range ticker.C{
		GetTokenHttp()
		log.Println("update the token,now token is:", token)
	}
}
func GetTokenHttp()  {
	req, err := http.NewRequest("GET",config.Config.TokenUrl, nil)
	if err != nil {
		log.Print(err)
	}
	q := req.URL.Query()
	q.Add("machineid", config.Config.MachineId)
	q.Add("password", config.Config.Password)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err!=nil{
		log.Print("error to update token")
		log.Print(err.Error())
	}
	resp_body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(resp_body))

	var tokenJson TokenJson
	err = json.Unmarshal(resp_body,&tokenJson)
	if err!=nil{
		log.Printf("error to parse json:")
		log.Printf(string(resp_body))
		log.Fatal(err.Error())
	}
	if tokenJson.Message=="success"{
		mutex.Lock()
		token = tokenJson.Token
		mutex.Unlock()
		log.Println("update the token,now token is: ", token)
	}
}

func GetToken() string  {
	mutex.Lock()
	nowToken := token
	mutex.Unlock()
	return nowToken
}

