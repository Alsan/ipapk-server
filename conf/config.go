package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var AppConfig *Config

type Config struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Proxy string `json:"proxy"`
}

func Init(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &AppConfig); err != nil {
		log.Fatal(err)
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

func (c *Config) ProxyURL() string {
	if c.Proxy == "" {
		return "http://" + c.Addr() + "/ipapk"
	}
	return c.Proxy + "/ipapk"
}
