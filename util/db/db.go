package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type DbConf struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DbName   string `yaml:"dbname" json:"dbname"`
}

var dbCli *sql.DB

func InitDb(conf *DbConf) error {
	if conf == nil {
		return errors.New("db config is nil")
	}
	var err error
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Password, conf.DbName)
	dbCli, err = sql.Open("postgres", connectStr)
	if err != nil {
		log.Println("fail to connect DB:" + err.Error())
		return err
	}
	return nil
}

func GetDbClient() *sql.DB {
	return dbCli
}
