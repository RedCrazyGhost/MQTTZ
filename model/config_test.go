package model

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestParseConfig(t *testing.T) {
	yamlStr := `
mqtt_configs:
- broker: 127.0.0.1
  port: 1883
  client_id: mqtt_client_1
  username: "emqx"
  password: "emqx2024"
  nickname: "MQTTZ_1" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key
  pub_configs:
  - enable_for: false
    interval: 1s
    source_path: ./json/data.json
    source_type: json
  - enable_for: false
    interval: 1s
    source_path: ./yaml/data.yaml
    source_type: yaml
  - enable_for: false
    interval: 1s
    source_type: conf
    source_data:
    - topic: 1
    - topic: 2
- broker: 127.0.0.1
  port: 1883
  client_id: mqtt_client_2
  username: "emqx"
  password: "emqx2024"
  nickname: "MQTTZ_2" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key
- broker: 127.0.0.1
  port: 1883
  client_id: mqtt_client_3
  username: "emqx"
  password: "emqx2024"
  nickname: "MQTTZ_3" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key
`
	var config Config
	err := yaml.Unmarshal([]byte(yamlStr), &config)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", config)
}
