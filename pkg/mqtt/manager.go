package mqtt

import (
	"errors"
	"sync"

	"MQTTZ/model"
	"MQTTZ/pkg/logger"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewMQTTClientManager)

type ClientManager struct {
	clientMaps sync.Map // map[string]*Client

	forwardRulesMap map[string]map[string][]model.ForwardRule // clientID -> topic -> forwardRules
}

func NewMQTTClientManager(conf *model.Config) (*ClientManager, error) {
	manager := &ClientManager{
		clientMaps:      sync.Map{},
		forwardRulesMap: make(map[string]map[string][]model.ForwardRule),
	}

	for i := 0; i < len(conf.MQTTConfigs); i++ {
		config := conf.MQTTConfigs[i]
		client, err := NewMQTTClient(&config)
		if err != nil {
			return nil, err
		}

		forwardRules, ok := manager.forwardRulesMap[client.id]
		if !ok {
			forwardRules = make(map[string][]model.ForwardRule)
		}
		for _, subConfig := range config.SubConfigs {
			if subConfig.Topic != "" {
				rules, ok := forwardRules[subConfig.Topic]
				if !ok {
					rules = make([]model.ForwardRule, 0)
				}
				rules = append(rules, subConfig.ForwardRules...)
				forwardRules[subConfig.Topic] = rules
			}
			for _, topic := range subConfig.Topics {
				rules, ok := forwardRules[topic]
				if !ok {
					rules = make([]model.ForwardRule, 0)
				}
				rules = append(rules, subConfig.ForwardRules...)
				forwardRules[topic] = rules
			}
		}
		manager.forwardRulesMap[client.id] = forwardRules

		manager.clientMaps.Store(client.id, client)
	}

	logger.Info("forward rules map", zap.Any("forward_rules_map", manager.forwardRulesMap))
	return manager, nil
}

func (m *ClientManager) Start() {
	m.clientMaps.Range(func(_, value any) bool {
		client := value.(*Client)
		go client.Run()
		go func() {
			if err := m.MQTTClientForwardData(client.id); err != nil {
				logger.Error("forward data failed",
					zap.String("client_id", client.id),
					zap.Error(err),
				)
			}
		}()
		logger.Info("start mqtt client", zap.String("client_id", client.id))
		return true
	})

}

func (m *ClientManager) GetMQTTClient(key string) *Client {
	client, ok := m.clientMaps.Load(key)
	if !ok {
		return nil
	}
	return client.(*Client)
}

func (m *ClientManager) GetMQTTClientInputDataChan(key string) chan<- model.MQTTDataProtocol {
	client, ok := m.clientMaps.Load(key)
	if !ok {
		return nil
	}
	return client.(*Client).pubDataCh
}

func (m *ClientManager) GetMQTTClientOutputDataChan(key string) <-chan model.MQTTDataProtocol {
	client, ok := m.clientMaps.Load(key)
	if !ok {
		return nil
	}
	return client.(*Client).subDataCh
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

// 转发MQTT数据
func (m *ClientManager) MQTTClientForwardData(fromClient string) error {
	fromClientChan := m.GetMQTTClientOutputDataChan(fromClient)
	if fromClientChan == nil {
		return errors.New("mqtt client is not found")
	}

	for data := range fromClientChan {
		ruleList, ok := m.forwardRulesMap[fromClient][data.GetTopic()]
		if !ok {
			continue
		}
		for _, rule := range ruleList {
			d := data.(model.MQTTData)
			d.Topic = data.(model.MQTTData).Topic + "/test"
			_ = m.MQTTClientPub(rule.ToClient, d)
		}
	}

	return nil
}
