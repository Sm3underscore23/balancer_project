package model

import (
	"net/http/httputil"
	"sync"
	"sync/atomic"
)

type BackendServer struct {
	IsOnline atomic.Bool
	Prx      *httputil.ReverseProxy
	BckndUrl string
	Method   string
	HelthUrl string
}

type BackendPool struct {
}

func NewBackendPool(backend) *BackendPool {
	backendList := make([]*model.BackendServer, len(cfg.BackendList))
	for i, b := range cfg.BackendList {
		if b.Config.Health.URL == "" {
			b.Config.Health.URL = "/health"
		}
		if b.Config.Health.Method == "" {
			b.Config.Health.Method = "GET"
		}
		backendList[i] = &model.BackendServer{
			BckndUrl: b.BackendURL,
			Method:   b.Config.Health.Method,
			HelthUrl: b.BackendURL + b.Config.Health.URL,
		}
	}
	return &BackendPool{Pool: backendList}
}
