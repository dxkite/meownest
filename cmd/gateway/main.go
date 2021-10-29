package main

import (
	"context"
	"dxkite.cn/gateway/config"
	"dxkite.cn/gateway/route"
	"dxkite.cn/gateway/server"
	"dxkite.cn/gateway/util"
	"dxkite.cn/log"
	"flag"
	"net"
	"os"
)

func init() {
	log.SetOutput(log.NewColorWriter())
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func NewGateway(ctx context.Context, p string) {
	cfg := config.NewConfig()
	if err := cfg.LoadFromFile(p); err != nil {
		log.Error(err)
	}
	util.SetLogConfig(ctx, cfg.LogConfig.Level, cfg.LogConfig.Path)
	log.Println("load config success")
	r := route.NewRoute()
	r.Load(cfg.Routes)
	s := server.NewServer(cfg, r)
	if cfg.HotLoad > 0 {
		cfg.SetLoadTime(cfg.HotLoad)
		cfg.OnChange(func(c interface{}) {
			cc := c.(*config.Config)
			cfg.SetLoadTime(cc.HotLoad)
			r.ClearAll()
			r.Load(cc.Routes)
			s.ApplyHeaderFilter(cc.HttpAllowHeader)
			s.ApplyCorsConfig(cc.Cors)
		})
		cfg.HotLoadIfModify(p)
	}
	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error(err)
	}
	log.Println("gateway start at", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Error(err)
	}
}

func main() {
	ctx, _ := context.WithCancel(context.Background())
	conf := flag.String("conf", "./config.yml", "the config file")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	NewGateway(ctx, *conf)
}
