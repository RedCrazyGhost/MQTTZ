package model

import (
	"errors"

	"MQTTZ/utils"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LogConfig struct {
	EnableDebug bool   `yaml:"enable_debug"`
	OutputLevel string `yaml:"level"`
}

type MQTTConfig struct {
	Broker   string `yaml:"broker"`
	Port     int    `yaml:"port"`
	ClientID string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Nickname string `yaml:"nickname"`

	PubConfigs []PubConfig `yaml:"pub_configs"`
	SubConfigs []SubConfig `yaml:"sub_configs"`
}

type Config struct {
	Log    LogConfig    `yaml:"log"`
	Server ServerConfig `yaml:"server"`

	MQTTConfigs []MQTTConfig `yaml:"mqtt_configs"`
}

type PubConfig struct {
	EnableFor  bool       `yaml:"enable_for,omitempty"`  // 是否循环
	Interval   string     `yaml:"interval,omitempty"`    // 每个 topic 数据发送的时间间隔
	SourceType SourceType `yaml:"source_type"`           // 数据源类型
	SourcePath string     `yaml:"source_path,omitempty"` // 数据源位置
	SourceData []MQTTData `yaml:"source_data,omitempty"` // 数据源数据
}

type SourceType string

const (
	SourceTypeConf SourceType = "conf"
	SourceTypeJSON SourceType = "json"
	SourceTypeYAML SourceType = "yaml"
)

func (c *PubConfig) ParseData() error {
	switch c.SourceType {
	case SourceTypeConf:
		if len(c.SourceData) == 0 {
			return errors.New("请输入发送数据！发送数据不能为空")
		}
		return nil
	case SourceTypeJSON:
		err := utils.LoadJSONFile(c.SourcePath, &c.SourceData)
		if err != nil {
			return err
		}
	case SourceTypeYAML:
		err := utils.LoadYAMLFile(c.SourcePath, &c.SourceData)
		if err != nil {
			return err
		}
	default:
		return errors.New("未知数据类型，请检查类型")
	}
	return nil
}

type SubConfig struct {
	Topic  string   `yaml:"topic,omitempty"`
	Topics []string `yaml:"topics,omitempty"`
	Qos    byte     `yaml:"qos,omitempty"`
}
