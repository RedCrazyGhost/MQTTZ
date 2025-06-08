package mqtt

import (
	"cmp"
	"fmt"
	"sync"
	"time"

	"MQTTZ/model"
	"MQTTZ/pkg/logger"
	"MQTTZ/utils/color"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Client struct {
	id string
	c  mqtt.Client
	sync.WaitGroup

	pubConf []model.PubConfig
	subConf []model.SubConfig

	pubDataCh chan model.MQTTDataProtocol
	subDataCh chan model.MQTTDataProtocol
}

func NewMQTTClient(conf *model.MQTTConfig) (*Client, error) {
	client := &Client{
		id: cmp.Or(conf.Nickname, fmt.Sprintf("%s@%s", conf.ClientID, conf.Broker)),

		pubConf:   conf.PubConfigs,
		subConf:   conf.SubConfigs,
		pubDataCh: make(chan model.MQTTDataProtocol, 999),
		subDataCh: make(chan model.MQTTDataProtocol, 999),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.Broker, conf.Port))
	opts.SetClientID(conf.ClientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logger.Info("MQTT client connected",
			zap.String("client_id", client.id),
		)
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logger.Error("MQTT client connection lost",
			zap.String("client_id", client.id),
			zap.Error(err),
		)
	})

	client.c = mqtt.NewClient(opts)

	return client, nil
}

func (m *Client) Run() {
	if token := m.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	m.sub()
	m.Add(1)
	go m.pub()

	for _, conf := range m.pubConf {
		m.Add(1)
		go func() {
			duration, err := time.ParseDuration(conf.Interval)
			if err != nil {
				duration = time.Second
			}
			defer m.Done()
			isFrist := true
			for isFrist || conf.EnableFor {
				isFrist = false
				for _, data := range conf.SourceData {
					m.Pub(data)
					time.Sleep(duration)
				}
			}
		}()
	}

	m.Wait()
}

func (m *Client) sub() {
	subMap := make(map[string]byte)
	for _, conf := range m.subConf {
		if conf.Topic != "" {
			subMap[conf.Topic] = conf.Qos
		}
		for _, topic := range conf.Topics {
			subMap[topic] = conf.Qos
		}
	}
	m.c.SubscribeMultiple(subMap, func(_ mqtt.Client, message mqtt.Message) {
		logger.Info(
			color.Theme.Sub.Text("sub"),
			zap.String("client_id", m.id),
			zap.String("topic", message.Topic()),
			zap.ByteString("payload", message.Payload()),
		)
		m.subDataCh <- model.MQTTData{
			Topic:   message.Topic(),
			QoS:     message.Qos(),
			Retain:  message.Retained(),
			Payload: message.Payload(),
		}
	})
}

func (m *Client) pub() {
	defer m.Done()
	for data := range m.pubDataCh {
		_ = m.c.Publish(data.GetTopic(), data.GetQoS(), data.GetRetain(), data.GetPayload())
	}
}

func (m *Client) Pub(data any) {
	mqttData, ok := data.(model.MQTTDataProtocol)
	if !ok {
		return
	}
	logger.Info(
		color.Theme.Pub.Text("pub"),
		zap.String("client_id", m.id),
		zap.String("topic", mqttData.GetTopic()),
		zap.ByteString("payload", mqttData.GetPayload()),
	)
	m.pubDataCh <- mqttData
}

func (m *Client) Sub() <-chan model.MQTTDataProtocol {
	return m.subDataCh
}
