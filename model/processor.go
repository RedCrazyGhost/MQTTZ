package model

type MQTTDataProcessorProtocol interface {
	GetRule() string
	GetProcessorType() ProcessorType
}

type ProcessorTargetType string

const (
	ProcessorTargetTypeTopic   ProcessorTargetType = "topic"
	ProcessorTargetTypeQos     ProcessorTargetType = "qos"
	ProcessorTargetTypeRetain  ProcessorTargetType = "retain"
	ProcessorTargetTypePayload ProcessorTargetType = "payload"
)

type ProcessorType string

const (
	ProcessorTypeInterceptor ProcessorType = "interceptor" // 拦截器
	ProcessorTypeFilter      ProcessorType = "filter"      // 过滤器
	ProcessorTypeForwarder   ProcessorType = "forwarder"   // 转发器
	ProcessorTypeExtractor   ProcessorType = "extractor"   // 数据提取
	ProcessorTypeGenerator   ProcessorType = "generator"   // 数据生成
)

type Processor struct {
	Type ProcessorType `yaml:"type" json:"type"`
	Rule string        `yaml:"rule" json:"rule"`
}

func (p Processor) GetRule() string {
	return p.Rule
}
func (p Processor) GetProcessorType() ProcessorType {
	return p.Type
}
