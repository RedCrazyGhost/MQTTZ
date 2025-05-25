package mqtt

import (
	"errors"
	"sync"

	"MQTTZ/model"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMQTTClientManager)

type ClientManager struct {
	clientMaps sync.Map // map[string]*MQTTClient
}

func NewMQTTClientManager(conf *model.Config) (*ClientManager, error) {
	manager := &ClientManager{
		clientMaps: sync.Map{},
	}

	for _, config := range conf.MQTTConfigs {
		client, err := NewMQTTClient(config)
		if err != nil {
			return nil, err
		}

		manager.clientMaps.Store(client.id, client)
	}

	return manager, nil
}

func (m *ClientManager) Start() {
	m.clientMaps.Range(func(_, value any) bool {
		go value.(*MQTTClient).Run()
		return true
	})
}

func (m *ClientManager) GetMQTTClient(key string) *MQTTClient {
	client, ok := m.clientMaps.Load(key)
	if !ok {
		return nil
	}
	return client.(*MQTTClient)
}

func (m *ClientManager) GetMQTTClientInputDataChan(key string) chan<- model.MQTTDataProtocol {
	client, ok := m.clientMaps.Load(key)
	if !ok {
		return nil
	}
	return client.(*MQTTClient).pubDataCh
}

func (m *ClientManager) MQTTClientPub(key string, data any) error {
	mqttData, ok := data.(model.MQTTDataProtocol)
	if !ok {
		return errors.New("data type is not model.MQTTDataProtocol")
	}
	client := m.GetMQTTClient(key)
	if client == nil {
		return errors.New("mqtt client is not found")
	}
	client.Pub(mqttData)
	return nil
}
