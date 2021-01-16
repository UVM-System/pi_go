package config

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	Config conf
	once   sync.Once
)

type capConfig struct {
	VideoId int    `yaml:"videoid"`
	Prefix  string `yaml:"prefix"`
}

type conf struct {
	VideoId    int         `yaml:"videoid"`
	SerialPort string      `yaml:"port"`
	Baudrate   int         `yaml:"baudrate"`
	DetectUrl  string      `yaml:"detecturl"`
	Pi2Url     string      `yaml:"pi2url"`
	CapConfigs []capConfig `yaml:"caps"`
	MachineId  string      `yaml:"machineid"`
	Password   string      `yaml:"password"`
	DoorPin    int         `yaml:"doorpin"`
}

func (c *conf) getConf() {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func init() {
	once.Do(readConfig)
}

func readConfig() {
	Config.getConf()
}
