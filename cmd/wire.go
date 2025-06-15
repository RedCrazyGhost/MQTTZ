//go:build wireinject
// +build wireinject

package main

import (
	"MQTTZ/pkg/conf"
	"MQTTZ/pkg/mqtt"
	"MQTTZ/pkg/processor"

	"github.com/google/wire"
)

func InitializeMQTTZ() (*MQTTZ, error) {
	panic(wire.Build(
		conf.ProviderSet,
		mqtt.ProviderSet,
		processor.ProviderSet,
		NewMQTTZ,
	))
	return &MQTTZ{}, nil
}
