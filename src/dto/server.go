package dto

import (
	"time"
)

// 服务
type Server struct {
	Id          string        `json:"id"`
	ServerName  []*ServerName `json:"server_name"`
	Collections []*Collection `json:"collections"`
	Endpoint    []*Endpoint   `json:"endpoint,omitempty"`
}

// 域名
type ServerName struct {
	Id          string       `json:"id"`
	Name        string       `json:"string"`      // 域名
	Protocol    string       `json:"protocol"`    // 协议
	Certificate *Certificate `json:"certificate"` // 证书
}

// SSL证书
type Certificate struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      []string  `json:"domain"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Key         string    `json:"key"`
	Certificate string    `json:"certificate"`
}

// 路由信息
type Route struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Method   []string    `json:"method"`
	Path     string      `json:"path"`
	Endpoint []*Endpoint `json:"endpoints,omitempty"`
}

// 路由组
type Collection struct {
	Id            string        `json:"id"`
	ParentId      string        `json:"parent_id,omitempty"` // 父级ID
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Authorization string        `json:"authorization"` // 权限校验
	Routes        []*Route      `json:"routes"`
	Collections   []*Collection `json:"collections,omitempty"`
	Endpoint      []*Endpoint   `json:"endpoints,omitempty"`
}

// 后端配置
type Endpoint struct {
	Id string `json:"id"`
	// 后端名
	Name string `json:"name"`
	// 服务类型
	Type string `json:"type"`
	// 转发正则
	ForwardRegex string `json:"forward_regex"`
	// 转发配置
	ForwardReplace string `json:"forward_replace"`
	// 请求头转发配置
	ForwardHeader []*ForwardHeaderOption `json:"forward_header"`
	// 匹配规则
	MatchFilter []*MatchOption `json:"match_filter"`
}

type ForwardHeaderOption struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MatchOption struct {
	Source string `json:"source"` // 匹配源
	Name   string `json:"name"`   // 匹配值
	Type   string `json:"type"`   // 匹配方式
	Value  string `json:"value"`  // 匹配内容
}
