package model

// 入口配置信息
type ServerName struct {
	Base
	Name          string `json:"string"`         // 域名
	Protocol      string `json:"protocol"`       // 协议
	CertificateId string `json:"certificate_id"` // 证书ID
}