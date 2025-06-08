package conf

import (
	"cmp"
	"encoding/json"
	"errors"
	"flag"

	"MQTTZ/model"
	"MQTTZ/pkg/logger"
	"MQTTZ/utils"

	"github.com/google/wire"
	"go.uber.org/zap"
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

	if err := logger.Init(&config.Log); err != nil {
		return nil, err
	}

	config.Server.Port = cmp.Or(config.Server.Port, _DefaultServerPort)
	config.Log.EnableDebug = cmp.Or(config.Log.EnableDebug, _DefaultDebug)

	for i := 0; i < len(config.MQTTConfigs); i++ {
		mqttConfig := config.MQTTConfigs[i]
		for j := 0; j < len(mqttConfig.PubConfigs); j++ {
			err = ParseData(&mqttConfig.PubConfigs[j])
			if err != nil {
				return nil, err
			}
		}
	}

	logger.Info("初始化：",
		zap.Any("config", config),
	)
	return config, nil
}

func ParseData(c *model.PubConfig) error {
	switch c.SourceType {
	case model.SourceTypeConf:
		if len(c.SourceData) == 0 {
			return errors.New("请输入发送数据！发送数据不能为空")
		}
	case model.SourceTypeJSON:
		err := utils.LoadJSONFile(c.SourcePath, &c.SourceData)
		if err != nil {
			return err
		}
	case model.SourceTypeYAML:
		err := utils.LoadYAMLFile(c.SourcePath, &c.SourceData)
		if err != nil {
			return err
		}
	default:
		return errors.New("未知数据类型，请检查类型")
	}

	var err error
	for i := 0; i < len(c.SourceData); i++ {
		switch d := c.SourceData[i].(type) {
		case map[string]any:
			var data any
			dataJSON, _ := json.Marshal(d)

			data = new(model.MockMQTTData)
			err = json.Unmarshal(dataJSON, &data)
			if err != nil {
				logger.Warn("json unmarshal mock mqtt data error", zap.Error(err))
			} else {
				c.SourceData[i] = data
				logger.Info("mock mqtt data", zap.Any("data", c.SourceData[i]))
				continue
			}

			data = new(model.MQTTData)
			err = json.Unmarshal(dataJSON, &data)
			if err != nil {
				logger.Error("json unmarshal mqtt data error", zap.Error(err))
				return err
			}

			c.SourceData[i] = data
			logger.Info("mqtt data", zap.Any("data", c.SourceData[i]))
		default:
			return errors.New("数据格式不合规")
		}
	}

	return nil
}
