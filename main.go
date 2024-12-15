package main

import (
	"encoding/json"
	"fmt"
	"mqtt/utils"
	"os"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	// makefile build auto generate
	AUTHOR  string
	VERSION string
)

var LOGO = strings.Join([]string{
	"  __  __  \033[1;38;2;98;191;140m___\033[0m\033[1;38;2;160;195;250m _____ _____\033[0m\033[1;38;2;255;0;0m _____\033[0m",
	" |  \\/  |\033[1;38;2;98;191;140m/ _ \\\033[0m\033[1;38;2;160;195;250m_   _|_   _|\033[0m\033[1;38;2;255;0;0m__  /\033[0m",
	" | |\\/| |\033[1;38;2;98;191;140m | | |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m   / /\033[0m",
	" | |  | | \033[1;38;2;98;191;140m|_| |\033[0m\033[1;38;2;160;195;250m| |   | |\033[0m\033[1;38;2;255;0;0m  / /_\033[0m",
	" |_|  |_|\033[1;38;2;98;191;140m\\__\\_\\\033[0m\033[1;38;2;160;195;250m|_|   |_|\033[0m\033[1;38;2;255;0;0m /____|\033[0m",
}, "\n") + "\n"

var (
	configFile = "config.json"
	conf       = &Config{}
	waitGroup  = sync.WaitGroup{}
	outputFile *os.File
	mqttClient mqtt.Client
)

type MQTTConfig struct {
	Broker   string `json:"broker"`
	Port     int    `json:"port"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	MQTTConfig
	InputConfigs []*InputConfig `json:"input_configs"`
	OutPutConfig OutputConfig   `json:"output_config"`
}

type InputConfig struct {
	IsFor    bool       `json:"is_for"`   // 是否循环
	Interval string     `json:"interval"` // 每个 topic 数据发送的时间间隔
	Source   string     `json:"source"`   // 数据源，没有使用自身数据，有则加载json
	Data     []MQTTData `json:"mqtt_data"`
}

type OutputConfig struct {
	SubTopics      []string `json:"sub_topics"`
	OutputFileName string   `json:"output_file_name"`
}

type MQTTData struct {
	Topic string `json:"topic"`
	Data  any    `json:"data"`
}

func main() {
	fmt.Println(LOGO)
	fmt.Printf("Author: %s\n", AUTHOR)
	fmt.Printf("Version: %s\n", VERSION)
	loadConfig()
	fmt.Printf("%+v\n", conf)
	NewMqttClient()
	subData()
	for _, config := range conf.InputConfigs {
		waitGroup.Add(1)
		go pubData(*config)
	}
	waitGroup.Wait()
	_, err := outputFile.WriteString("]")
	if err != nil {
		panic(err)
	}
	fmt.Println("finished")
}

func NewMqttClient() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.Broker, conf.Port))
	opts.SetClientID(conf.ClientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
		os.Exit(0)
	}
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// loadConfig 使用 config.json
func loadConfig() {
	utils.LoadJSONFile(configFile, &conf)

	for _, config := range conf.InputConfigs {
		if config.Source != "" {
			utils.LoadJSONFile(config.Source, &config.Data)
		}
	}
}

func subData() {
	filters := make(map[string]byte)
	for _, topic := range conf.OutPutConfig.SubTopics {
		filters[topic] = 0
	}

	outputFileName := conf.OutPutConfig.OutputFileName
	if outputFileName == "" {
		outputFileName = time.Now().Format(time.RFC3339)
	}

	var err error
	outputFile, err = os.OpenFile(outputFileName+".json", os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	_, err = outputFile.WriteString("[")
	if err != nil {
		panic(err)
	}
	isSubFrist := true

	mqttClient.SubscribeMultiple(filters, func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Printf("%s %s topic:%s data:%v\n",
			time.Now().Format(time.RFC3339),
			utils.UseColor("sub", 255, 255, 0),
			msg.Topic(),
			string(msg.Payload()))

		var d any
		_ = json.Unmarshal(msg.Payload(), &d)
		data := MQTTData{
			Topic: msg.Topic(),
			Data:  d,
		}

		marshal, _ := json.Marshal(data)
		if !isSubFrist {
			_, err := outputFile.WriteString(",")
			if err != nil {
				panic(err)
			}
		}
		_, err = outputFile.Write(marshal)
		if err != nil {
			panic(err)
		}
		isSubFrist = false
	})
}

func pubData(conf InputConfig) {
	isFrist := true
	for isFrist || conf.IsFor {
		isFrist = false
		for _, data := range conf.Data {
			fmt.Printf("%s %s topic:%s data:%v\n",
				time.Now().Format(time.RFC3339),
				utils.UseColor("pub", 0, 0, 255),
				data.Topic,
				data.Data)

			duration, err := time.ParseDuration(conf.Interval)
			if err != nil {
				duration = time.Second
			}
			marshal, err := json.Marshal(data.Data)
			if err != nil {
				continue
			}
			publish := mqttClient.Publish(data.Topic, 0, false, marshal)
			publish.Wait()
			time.Sleep(duration)
		}
	}
	waitGroup.Done()
}
