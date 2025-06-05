package main

import (
	"fmt"
	"strings"

	"MQTTZ/model"
	"MQTTZ/pkg/mqtt"
)

// makefile build auto generate
var (
	AUTHOR   string
	VERSION  string
	REPO_URL string
)

var LOGO = strings.Join([]string{
	"  __  __  \033[1;38;2;98;191;140m___\033[0m\033[1;38;2;160;195;250m _____ _____\033[0m\033[1;38;2;255;0;0m _____\033[0m",
	" |  \\/  |\033[1;38;2;98;191;140m/ _ \\\033[0m\033[1;38;2;160;195;250m_   _|_   _|\033[0m\033[1;38;2;255;0;0m__  /\033[0m",
	" | |\\/| |\033[1;38;2;98;191;140m | | |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m   / /\033[0m",
	" | |  | | \033[1;38;2;98;191;140m|_| |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m  / /_\033[0m",
	" |_|  |_|\033[1;38;2;98;191;140m\\__\\_\\\033[0m\033[1;38;2;160;195;250m|_|   |_|\033[0m\033[1;38;2;255;0;0m /____|\033[0m",
	"Author: " + AUTHOR,
	"Version: " + VERSION,
	"Repo Url: " + REPO_URL,
}, "\n") + "\n"

type MQTTZ struct {
	config        *model.Config
	clientManager *mqtt.ClientManager
}

func main() {
	fmt.Println(LOGO)

	mqttz, err := InitializeMQTTZ()
	if err != nil {
		panic(err)
	}
	mqttz.clientManager.Start()

	go func() {
		targetClient := mqttz.clientManager.GetMQTTClient("MQTTZ_2")
		for protocol := range mqttz.clientManager.GetMQTTClient("MQTTZ_1").Sub() {
			targetClient.Pub(model.MQTTData{
				Topic:   "remove/" + protocol.GetTopic(),
				Payload: protocol.GetPayload(),
			})
		}
	}()

	select {}
}

func NewMQTTZ(config *model.Config, manager *mqtt.ClientManager) (*MQTTZ, error) {
	return &MQTTZ{
		config:        config,
		clientManager: manager,
	}, nil
}
