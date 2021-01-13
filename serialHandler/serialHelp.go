package serialHandler
import (
	"github.com/tarm/serial"
	"log"
)
var (
	noCareBytes = []byte{'0','1','\n','\r'}
)
func Serial_init(name string,baud int) (*serial.Port,error) {
	c:=&serial.Config{
		Name:name,
		Baud:baud,
	}
	s,err:=serial.OpenPort(c)
	return s,err
}

func Start_listen(port *serial.Port,signal chan string)  {
	buf:=make([]byte,128)
	for {
		n,err:=port.Read(buf)
		if err!=nil{
			log.Fatal(err)
		}
		newBytes :=RemoveNoCareSet(buf[:n])
		if len(newBytes)!=0{
			str :=string(newBytes)
			signal<-str
		} else {
			continue
		}

	}
}
func Contain(array []string,obj string) bool {
	for i:=0;i<len(array) ;i++  {
		if(obj==array[i]){
			return true
		}
	}
	return false
}

func RemoveNoCareSet(ori []byte) []byte {
	var newByte []byte
	for i:=0;i<len(ori);i++{
		if !ContainByte(noCareBytes,ori[i]){
			newByte = append(newByte,ori[i])
		}
	}
	return newByte
}

func ContainByte(array []byte,obj byte) bool {
	for i:=0;i<len(array) ;i++  {
		if(obj==array[i]){
			return true
		}
	}
	return false
}

type SerialObserver interface {
	// 由关注者去完成他们的事情,比如开关门状态,比如拍摄照片,结算等
	RecvSerialStr(obstr string)
}

func SerialPortListen(name string,bau int,observers []SerialObserver)  {
	s,err :=Serial_init(name,bau)
	if err!=nil{
		panic(err.Error())
	}
	recvChan :=make(chan string)
	go Start_listen(s,recvChan)
	for{
		recvStr:=<-recvChan
		for i:=0;i<len(observers);i++{
			observers[i].RecvSerialStr(recvStr)
			//log.Println("lalala" + recvStr)
		}
	}
}

