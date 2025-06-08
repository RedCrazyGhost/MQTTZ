//go:build wireinject
// +build wireinject

package main

import (
	"MQTTZ/pkg/conf"
	"MQTTZ/pkg/mqtt"

	"github.com/google/wire"
)

func InitializeMQTTZ() (*MQTTZ, error) {
	panic(wire.Build(
		conf.ProviderSet,
		mqtt.ProviderSet,
		NewMQTTZ,
	))
	return &MQTTZ{}, nil
}
