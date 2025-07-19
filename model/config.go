package model

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LogConfig struct {
	EnableDebug bool   `yaml:"enable_debug"`
	Level       string `yaml:"level"`
	EnableColor bool   `yaml:"enable_color"`
	OutputFile  string `yaml:"output_file"`
	MaxSize     int    `yaml:"max_size"`    // 单个文件最大尺寸，单位：MB
	MaxBackups  int    `yaml:"max_backups"` // 最大保留文件数
	MaxAge      int    `yaml:"max_age"`     // 最大保留天数
	Compress    bool   `yaml:"compress"`    // 是否压缩
}

type MQTTConfig struct {
	Broker   string `yaml:"broker"`
	Port     int    `yaml:"port"`
	ClientID string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Nickname string `yaml:"nickname"`

	PubConfigs []PubConfig `yaml:"pub_configs"`
	SubConfigs []SubConfig `yaml:"sub_configs"`
}

type Config struct {
	Log    LogConfig    `yaml:"log"`
	Server ServerConfig `yaml:"server"`

	MQTTConfigs []MQTTConfig `yaml:"mqtt_configs"`
}

type PubConfig struct {
	EnableFor  bool       `yaml:"enable_for,omitempty"`  // 是否循环
	Interval   string     `yaml:"interval,omitempty"`    // 每个 topic 数据发送的时间间隔
	SourceType SourceType `yaml:"source_type"`           // 数据源类型
	SourcePath string     `yaml:"source_path,omitempty"` // 数据源位置
	SourceData []any      `yaml:"source_data,omitempty"` // 数据源数据
	IsStrong   bool       `yaml:"is_strong,omitempty"`   // 数据模型是否强匹配
}

type SourceType string

const (
	SourceTypeConf SourceType = "conf"
	SourceTypeJSON SourceType = "json"
	SourceTypeYAML SourceType = "yaml"
)

type SubConfig struct {
	Topic  string   `yaml:"topic,omitempty"`
	Topics []string `yaml:"topics,omitempty"`
	Qos    byte     `yaml:"qos,omitempty"`

	Processors   []Processor   `yaml:"processors,omitempty"`    // 前置处理器
	ForwardRules []ForwardRule `yaml:"forward_rules,omitempty"` // 转发规则
}

type ForwardRule struct {
	ToClient   string      `yaml:"to_client" json:"to_client"`
	Processors []Processor `yaml:"processors" json:"processors"`
}
