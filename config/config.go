package config

import (
	"encoding/json"
	"github.com/peakedshout/go-CFC/loger"
	"os"
)

type Config struct {
	ProxyServerHost ProxyServerHostConfig `json:"ProxyServerHost"`
	ProxyMethod     ProxyMethodConfig     `json:"ProxyMethod"`
	Setting         SettingConfig         `json:"Setting"`
}
type ProxyMethodConfig struct {
	Http struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Http"`
	Https struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Https"`
	Socks struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Socks"`
}
type ProxyServerHostConfig struct {
	ProxyServerAddr string `json:"ProxyServerAddr"`
	LinkProxyKey    string `json:"LinkProxyKey"`
}
type SettingConfig struct {
	ReLinkTime string `json:"ReLinkTime"`
	LogLevel   uint8  `json:"LogLevel"`
	LogStack   bool   `json:"LogStack"`
}

func ReadConfig(path string) *Config {
	b, err := os.ReadFile(path)
	if err != nil {
		loger.SetLogError(err)
	}
	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		loger.SetLogError(err)
	}
	return &config
}
