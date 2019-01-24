package util

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Config struct{
	User string `json:"user"`
	Password string `json:"password"`
	Dbname string `json:"dbname"`
	Sqlhost string `json:"sqlhost"`
	Sqlport string `json:"sqlport"`
	Webserver string `json:"webserver"`
	Webport string `json:"webport"`
}

var Cfg *Config
var Fullurl string

func LoadConfig()*Config{
	raw, err := ioutil.ReadFile("config/config.json")
	if err != nil{
		fmt.Println(err.Error())
	}
	configuration := Config{}
	err = json.Unmarshal(raw, &configuration)
	return &configuration
}

func init(){
	Cfg = LoadConfig()
	url := fmt.Sprintf("%s:%s/",Cfg.Webserver,Cfg.Webport)
	Fullurl = url
}
