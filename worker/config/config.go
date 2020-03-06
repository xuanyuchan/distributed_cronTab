package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdClientConfig
}

type EtcdClientConfig struct {
	EtcdEndPoints     []string `json:"etcdEndPoints"`
	EtcdClientTimeOut int      `json:"etcdClientTimeOut"`
}

var (
	G_Config *Config
)

func ParseConfig(fileName string) error {
	G_Config = &Config{}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, G_Config)
	if err != nil {
		return err
	}
	return nil
}
