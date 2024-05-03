package agent

import "net/http"

type RequestMatcher interface {
	MatchRequest(req *http.Request) bool
}

type RequestForwardHandler interface {
	HandleRequest(w http.ResponseWriter, req *http.Request)
}

type AuthorizeHandler interface {
	HandleAuthorizeCheck(w http.ResponseWriter, req *http.Request) bool
}

type ForwardHandler interface {
	RequestMatcher
	RequestForwardHandler
}

type Handler struct {
	items []ForwardHandler
}

func NewHandler() *Handler {
	h := new(Handler)
	h.items = []ForwardHandler{}
	return h
}

func (h *Handler) Add(item ForwardHandler) {
	h.items = append(h.items, item)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, item := range h.items {
		// 匹配请求
		if item.MatchRequest(req) {
			// 进行权限校验
			if auth, ok := item.(AuthorizeHandler); ok {
				if !auth.HandleAuthorizeCheck(w, req) {
					return
				}
			}
			// 校验通过
			item.HandleRequest(w, req)
			return
		}
	}

	// 无匹配路由
	http.NotFound(w, req)
}

type forwardItem struct {
	auth AuthorizeHandler
	RequestMatcher
	RequestForwardHandler
}

func NewForwardHandler(matcher RequestMatcher, forward RequestForwardHandler, auth AuthorizeHandler) ForwardHandler {
	return &forwardItem{RequestMatcher: matcher, RequestForwardHandler: forward, auth: auth}
}

func (item forwardItem) HandleAuthorizeCheck(w http.ResponseWriter, req *http.Request) bool {
	if item.auth != nil {
		return item.auth.HandleAuthorizeCheck(w, req)
	}
	return true
}
