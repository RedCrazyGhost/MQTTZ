package main

import (
	"fmt"
	"strings"

	"MQTTZ/model"
	"MQTTZ/pkg/logger"
	"MQTTZ/pkg/mqtt"

	"go.uber.org/zap"
)


// makefile build auto generate.
var (
	Author  string
	Version string
	RepoURL string
)

var LOGO = strings.Join([]string{
	"  __  __  \033[1;38;2;98;191;140m___\033[0m\033[1;38;2;160;195;250m _____ _____\033[0m\033[1;38;2;255;0;0m _____\033[0m",
	" |  \\/  |\033[1;38;2;98;191;140m/ _ \\\033[0m\033[1;38;2;160;195;250m_   _|_   _|\033[0m\033[1;38;2;255;0;0m__  /\033[0m",
	" | |\\/| |\033[1;38;2;98;191;140m | | |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m   / /\033[0m",
	" | |  | | \033[1;38;2;98;191;140m|_| |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m  / /_\033[0m",
	" |_|  |_|\033[1;38;2;98;191;140m\\__\\_\\\033[0m\033[1;38;2;160;195;250m|_|   |_|\033[0m\033[1;38;2;255;0;0m /____|\033[0m",
	" Author: " + Author,
	" Version: " + Version,
	" Repo Url: " + RepoURL,
}, "\n") + "\n"

type MQTTZ struct {
	config        *model.Config
	clientManager *mqtt.ClientManager
}

func main() {
	fmt.Println(LOGO)

	mqttz, err := InitializeMQTTZ()
	if err != nil {
		logger.Error("initialize mqttz error", zap.Error(err))
		panic(err)
	}
	logger.Info("initialize mqttz success")
	mqttz.clientManager.Start()
	select {}
}

func NewMQTTZ(config *model.Config, manager *mqtt.ClientManager) *MQTTZ {
	return &MQTTZ{
		config:        config,
		clientManager: manager,
	}
}
