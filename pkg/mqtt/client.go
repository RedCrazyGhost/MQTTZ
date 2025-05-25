package mqtt

import (
	"cmp"
	"fmt"
	"sync"
	"time"

	"MQTTZ/model"
	"MQTTZ/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	id string
	c  mqtt.Client
	sync.WaitGroup

	pubConf []model.PubConfig
	subConf []model.SubConfig

	pubDataCh chan model.MQTTDataProtocol
	subDataCh chan model.MQTTDataProtocol
}

func NewMQTTClient(conf model.MQTTConfig) (*MQTTClient, error) {
	client := &MQTTClient{
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
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Printf("Connect Client Name: %s\n", conf.ClientID)
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		fmt.Printf("Connect Lost Client Name: %s\n", conf.ClientID)
	})

	client.c = mqtt.NewClient(opts)

	return client, nil
}

func (m *MQTTClient) Run() {
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

func (m *MQTTClient) sub() {
	subMap := make(map[string]byte)
	for _, conf := range m.subConf {
		if conf.Topic != "" {
			subMap[conf.Topic] = conf.Qos
		}
		for _, topic := range conf.Topics {
			subMap[topic] = conf.Qos
		}
	}
	m.c.SubscribeMultiple(subMap, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("%s %s topic:%s data:%v\n",
			time.Now().Format(time.RFC3339),
			utils.UseColor("sub", 255, 255, 0),
			message.Topic(),
			string(message.Payload()))
		m.subDataCh <- model.MQTTData{
			Topic:   message.Topic(),
			QoS:     message.Qos(),
			Retain:  message.Retained(),
			Payload: message.Payload(),
		}
	})
}

func (m *MQTTClient) pub() {
	defer m.Done()
	for data := range m.pubDataCh {
		_ = m.c.Publish(data.GetTopic(), data.GetQoS(), data.GetRetain(), data.GetPayload())
	}
}

func (m *MQTTClient) Pub(data any) {
	mqttData, ok := data.(model.MQTTDataProtocol)
	if !ok {
		return
	}
	fmt.Printf("%s %s topic:%s data:%v\n",
		time.Now().Format(time.RFC3339),
		utils.UseColor("pub", 0, 0, 255),
		mqttData.GetTopic(),
		string(mqttData.GetPayload()))
	m.pubDataCh <- mqttData
}

func (m *MQTTClient) Sub() <-chan model.MQTTDataProtocol {
	return m.subDataCh
}
