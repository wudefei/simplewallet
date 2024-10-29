package config

import (
	"flag"
	"fmt"
	"os"
	"simplewallet/util/db"

	"gopkg.in/yaml.v2"
)

var Config Conf

type Conf struct {
	Env     string       `yaml:"env"`
	GinHost string       `yaml:"gin_host"`
	Db      db.DbConf    `yaml:"db"`
	Redis   db.RedisConf `yaml:"redis"`
}

var gConfigName string

func init() {
	fmt.Println("config init")
	flag.StringVar(&gConfigName, "conf", "./conf.yaml", "config name")
	flag.Parse()

	fmt.Println("config name ", gConfigName)

	ParseYaml(gConfigName, &Config)
}

func ParseYaml(file string, configRaw interface{}) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic("load config file error! reason:" + err.Error())
	}

	err = yaml.Unmarshal(content, configRaw)
	if err != nil {
		panic("Unmarshal config file error! reason:" + err.Error())
	}
}
