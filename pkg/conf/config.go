package conf

import (
	"cmp"
	"flag"
	"fmt"

	"MQTTZ/model"
	"MQTTZ/utils"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConfig)

const (
	_DefaultSConfigFile = "./conf/config.yaml"
	_DefaultServerPort  = 10514
	_DefaultDebug       = false
)

var configFile string

func NewConfig() (*model.Config, error) {
	config := &model.Config{}

	flag.StringVar(&configFile, "config", _DefaultSConfigFile, "运行配置文件地址")
	flag.BoolVar(&config.Log.EnableDebug, "debug", _DefaultDebug, "是否启发调试")
	flag.IntVar(&config.Server.Port, "port", _DefaultServerPort, "服务端口")
	flag.Parse()

	configFile = cmp.Or(configFile, _DefaultSConfigFile)

	err := utils.LoadYAMLFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	config.Server.Port = cmp.Or(config.Server.Port, _DefaultServerPort)
	config.Log.EnableDebug = cmp.Or(config.Log.EnableDebug, _DefaultDebug)

	for i := 0; i < len(config.MQTTConfigs); i++ {
		mqttConfig := config.MQTTConfigs[i]
		for j := 0; j < len(mqttConfig.PubConfigs); j++ {
			err = mqttConfig.PubConfigs[j].ParseData()
			if err != nil {
				return nil, err
			}
		}
	}

	fmt.Printf("%+v\n", config)
	return config, nil
}
