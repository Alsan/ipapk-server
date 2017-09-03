package conf

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
)

var AppConfig *Config

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
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
