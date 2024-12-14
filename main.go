package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
	InputConfigs []InputConfig `json:"input_configs"`
	OutPutConfig OutputConfig  `json:"output_config"`
}

type InputConfig struct {
	IsFor    bool       `json:"is_for"`   // 是否循环
	Interval string     `json:"interval"` // 每个 topic 数据发送的时间间隔
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
	loadConfigFile()
	fmt.Printf("%+v\n", conf)
	NewMqttClient()
	subData()
	for _, config := range conf.InputConfigs {
		waitGroup.Add(1)
		go pubData(config)
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

// loadConfigFile 使用 config.json
func loadConfigFile() {
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
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
		fmt.Printf("sub topic:%s data:%v\n", msg.Topic(), string(msg.Payload()))

		var d any
		_ = json.Unmarshal([]byte(msg.Payload()), &d)
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
	for ; isFrist || conf.IsFor; {
		isFrist = false
		for _, data := range conf.Data {
			fmt.Printf("pub topic:%s data:%v\n", data.Topic, data.Data)

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
