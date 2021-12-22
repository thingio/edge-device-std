package models

type (
	ProductFuncID     = string        // product functionality ID
	ProductPropertyID = ProductFuncID // product property's functionality ID
	ProductEventID    = ProductFuncID // product event's functionality ID
	ProductMethodID   = ProductFuncID // product method's functionality ID
)

const (
	DeviceDataMultiPropsID ProductFuncID = "*"
)

type Product struct {
	ID         string             `json:"id"`                    // 产品 ID
	Name       string             `json:"name"`                  // 产品名称
	Desc       string             `json:"desc"`                  // 产品描述
	Protocol   string             `json:"protocol"`              // 产品协议
	DataFormat string             `json:"data_format,omitempty"` // 数据格式
	Properties []*ProductProperty `json:"properties,omitempty"`  // 属性功能列表
	Events     []*ProductEvent    `json:"events,omitempty"`      // 事件功能列表
	Methods    []*ProductMethod   `json:"methods,omitempty"`     // 方法功能列表
	Topics     []*ProductTopic    `json:"topics,omitempty"`      // 各功能对应的消息主题
}

type ProductProperty struct {
	Id         ProductPropertyID `json:"id"`
	Name       string            `json:"name"`
	Desc       string            `json:"desc"`
	Interval   string            `json:"interval"`
	Unit       string            `json:"unit"`
	FieldType  string            `json:"field_type"`
	ReportMode string            `json:"report_mode"`
	Writeable  bool              `json:"writeable"`
	AuxProps   map[string]string `json:"aux_props"`
}

type ProductEvent struct {
	Id       ProductEventID    `json:"id"`
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Outs     []*ProductField   `json:"outs"`
	AuxProps map[string]string `json:"aux_props"`
}

type ProductMethod struct {
	Id       ProductMethodID   `json:"id"`
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Ins      []*ProductField   `json:"ins"`
	Outs     []*ProductField   `json:"outs"`
	AuxProps map[string]string `json:"aux_props"`
}

type ProductTopic struct {
	Topic   string `json:"topic"`
	OptType string `json:"opt_type"`
	Desc    string `json:"desc"`
}

type ProductField struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	FieldType string `json:"field_type"`
	Desc      string `json:"desc"`
}
