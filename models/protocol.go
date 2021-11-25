package models

type (
	ProtocolPropertyKey = string // common property owned by devices with the same protocol
)

type Protocol struct {
	ID           string      `json:"id"`   // 协议 ID
	Name         string      `json:"name"` // 协议名称
	Desc         string      `json:"desc"` // 协议描述
	Category     string      `json:"category"`
	Language     string      `json:"language"`
	SupportFuncs []string    `json:"support_funcs"`
	AuxProps     []*Property `json:"aux_props"`
	DeviceProps  []*Property `json:"device_props"`
}
