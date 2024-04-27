package server

import (
	"net/http"

	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewRoute(s service.Route) *Route {
	return &Route{s: s}
}

type Route struct {
	s service.Route
}

func (s *Route) Create(c *gin.Context) {
	var param service.CreateRouteParam

	if err := c.ShouldBind(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	Result(c, http.StatusCreated, rst)
}

func (s *Route) Get(c *gin.Context) {
	var param service.GetRouteParam

	if err := c.ShouldBindUri(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}
	Result(c, http.StatusOK, rst)
}

func (s *Route) List(c *gin.Context) {
	var param service.ListRouteParam

	if err := c.ShouldBindQuery(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.List(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	Result(c, http.StatusOK, rst)
}

func (s *Route) Update(c *gin.Context) {
	var param service.UpdateRouteParam
	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Update(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}
	Result(c, http.StatusOK, rst)
}

func (s *Route) Delete(c *gin.Context) {
	var param service.DeleteRouteParam

	if err := c.ShouldBindUri(&param); err != nil {
		ResultErrorBind(c, err)
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		ResultError(c, err)
		return
	}

	ResultEmpty(c, http.StatusOK)
}

func WithRoute(path string, server *Route) func(s *HttpServer) {
	return func(s *HttpServer) {
		group := s.engine.Group(path)
		{
			group.POST("", server.Create)
			group.GET("", server.List)
			group.GET("/:id", server.Get)
			group.POST("/:id", server.Update)
			group.DELETE("/:id", server.Delete)
		}
	}
}
