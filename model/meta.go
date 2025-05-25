package model

import "encoding/json"

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
	bytes, err := json.Marshal(d.Payload)
	if err != nil {
		return nil
	}
	return bytes
}
