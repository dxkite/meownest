package route

import (
	"dxkite.cn/gateway/config"
	"dxkite.cn/log"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Route struct {
	r  map[string]int
	re []*routeEntry
	bg []*RouteInfo
}

func NewRoute() *Route {
	return &Route{
		r:  map[string]int{},
		re: nil,
		bg: nil,
	}
}

type routeEntry struct {
	pattern string
	index   int
}

func (r routeEntry) String() string {
	return fmt.Sprintf("{%s:%d}", r.pattern, r.index)
}

type Backend interface {
	BackendType() string
}

type UrlBackend struct {
	Url *url.URL
}

func (u *UrlBackend) BackendType() string {
	return u.Url.Scheme
}

func NewUrlBackend(rawurl string) (Backend, error) {
	if u, err := url.Parse(rawurl); err != nil {
		log.Error("parse backend error", rawurl)
		return nil, err
	} else {
		host := u.Host
		port := ""
		if strings.LastIndex(host, ":") > 0 {
			host, port, _ = net.SplitHostPort(u.Host)
		}
		if len(port) == 0 {
			switch u.Scheme {
			case "http":
				port = "80"
			case "https":
				port = "443"
			default:
				return nil, errors.New("missing port")
			}
		}
		return &UrlBackend{Url: u}, nil
	}
}

type BackendGroup []Backend
type RouteInfo struct {
	Config config.Route
	Group  *BackendGroup
}

// 支持 http/https
// https会验证客户端
func NewBackendGroupFromUrl(backends []string) *BackendGroup {
	bg := BackendGroup{}
	for _, b := range backends {
		if bb, err := NewUrlBackend(b); err != nil {
			log.Error("create url backend error", b)
		} else {
			bg = append(bg, bb)
		}
	}
	return &bg
}

// 随机获取一个后端
func (b BackendGroup) Get() Backend {
	n := len(b)
	if n == 1 {
		return b[0]
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	idx := rand.Intn(len(b))
	return b[idx]
}

// 载入路由
func (r *Route) Load(routes []config.Route) {
	for i, route := range routes {
		idx := len(r.bg)
		r.re = append(r.re, &routeEntry{
			pattern: route.Pattern,
			index:   idx,
		})
		r.bg = append(r.bg, &RouteInfo{
			Group:  NewBackendGroupFromUrl(route.Backend),
			Config: routes[i],
		})
		r.r[route.Pattern] = idx
	}
	// 优先处理长前缀
	sort.Slice(r.re, func(i, j int) bool {
		return len(r.re[i].pattern) > len(r.re[j].pattern)
	})
}

// 清空路由
func (r *Route) ClearAll() {
	*r = *NewRoute()
}

// 匹配路由
func (r *Route) Match(pattern string) (string, *RouteInfo) {
	// 完整路由
	if idx, ok := r.r[pattern]; ok {
		return pattern, r.bg[idx]
	}
	// 前缀路由
	for _, e := range r.re {
		if strings.HasPrefix(pattern, e.pattern) {
			return e.pattern, r.bg[e.index]
		}
	}
	return "", nil
}
