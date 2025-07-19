package processor

import (
	"fmt"
	"strings"

	"MQTTZ/model"
	"MQTTZ/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewValidate)

var v *validator.Validate

func NewValidate() *validator.Validate {
	v = validator.New(validator.WithRequiredStructEnabled())
	return v
}

func Do(processor model.MQTTDataProcessorProtocol, data model.MQTTDataProtocol) bool {
	var err error
	targetType, ruleStr, err := uncodeRule(processor.GetRule())
	if err != nil {
		return false
	}
	switch targetType {
	case model.ProcessorTargetTypeTopic:
		err = v.Var(data.GetTopic(), ruleStr)
	case model.ProcessorTargetTypePayload:
		err = v.Var(data.GetPayload(), ruleStr)
	case model.ProcessorTargetTypeQos:
		err = v.Var(data.GetQoS(), ruleStr)
	case model.ProcessorTargetTypeRetain:
		err = v.Var(data.GetRetain(), ruleStr)
	default:
		return false
	}
	logger.Warn("processor rule", zap.Any("rule", processor.GetRule()), zap.Any("data", data), zap.Any("err", err))

	switch processor.GetProcessorType() {
	case model.ProcessorTypeInterceptor:
		return err == nil
	case model.ProcessorTypeFilter:
		return err != nil
	}
	return true
}

func uncodeRule(rule string) (model.ProcessorTargetType, string, error) {
	index := strings.Index(rule, ":")
	if index == -1 {
		return "", "", fmt.Errorf("rule %s is invalid", rule)
	}
	return model.ProcessorTargetType(rule[:index]), rule[index+1:], nil
}
