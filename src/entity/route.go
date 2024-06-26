package entity

import (
	"time"

	"dxkite.cn/meownest/src/enum"
	"dxkite.cn/meownest/src/value"
)

type Route struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Method      []string `json:"method" gorm:"serializer:json"`
	//路径
	Path string `json:"path"`
	// 路径类型
	PathType enum.RoutePathType `json:"path_type"`
	// 匹配规则
	MatchOptions []*value.MatchOption `json:"match_options" gorm:"serializer:json"`
	// 路径重写
	PathRewrite *value.PathRewrite `json:"path_rewrite" gorm:"serializer:json"`
	// 数据编辑
	ModifyOptions []*value.ModifyOption `json:"modify_options" gorm:"serializer:json"`
	// 所属集合ID
	CollectionId uint64 `gorm:"index"`
	// 权限配置ID
	AuthorizeId uint64 `gorm:"index"`
	// 后端服务ID
	EndpointId uint64 `gorm:"index"`
	// 路由状态
	Status enum.RouteStatus
}

func NewRoute() *Route {
	return new(Route)
}
