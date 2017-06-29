package config

import (
	"github.com/go-ini/ini"
	"os"
	"xmn/core/logger"
)

func Read(section, keys string) string {
	path := os.Getenv("GOPATH") + "/src/xmn/core/config/config.ini"
	cnf, err := ini.Load(path)
	if err != nil {
		logger.Log("ini 配置文件 读取错误")
		os.Exit(0)
	}

	section_name := section
	key_name := keys

	val := cnf.Section(section_name).Key(key_name).String()

	if val == "" {
		logger.Log(key_name + " 不在 ini 配置项里")
	}
	return val
}