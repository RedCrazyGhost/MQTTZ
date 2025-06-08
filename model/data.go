package model

import (
	"github.com/brianvoe/gofakeit/v7"
	json "github.com/bytedance/sonic"
)

type MQTTDataProtocol interface {
	GetTopic() string
	GetRetain() bool
	GetQoS() byte
	GetPayload() []byte
}

type MQTTData struct {
	Topic   string `json:"topic" yaml:"topic"`
	Retain  bool   `json:"retain" yaml:"retain"`
	QoS     byte   `json:"qos" yaml:"qos"`
	Payload any    `json:"payload" yaml:"payload"`
}

func (d MQTTData) GetTopic() string {
	return d.Topic
}

func (d MQTTData) GetRetain() bool {
	return d.Retain
}

func (d MQTTData) GetQoS() byte {
	return d.QoS
}

func (d MQTTData) GetPayload() []byte {
	switch payload := d.Payload.(type) {
	case []byte:
		return payload
	default:
		bytes, err := json.Marshal(d.Payload)
		if err != nil {
			return nil
		}
		return bytes
	}
}

type MockMQTTData struct {
	MockData any `json:"mock_data" yaml:"mock_data"`
	MQTTData
}

func (d MockMQTTData) GetPayload() []byte {
	if d.MockData != nil {
		template, err := gofakeit.Template(d.Payload.(string), &gofakeit.TemplateOptions{Data: d.MockData})
		if err != nil {
			return nil
		}
		return []byte(template)
	}

	return d.MQTTData.GetPayload()
}
