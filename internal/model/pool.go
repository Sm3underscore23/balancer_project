package model

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type BackendServerSettings struct {
	BckndUrl string
	Method   string
	HelthUrl string
}

type BackendServer struct {
	isOnline atomic.Bool
	Prx      *httputil.ReverseProxy
	BackendServerSettings
}

func (b *BackendServer) Load() bool {
	return b.isOnline.Load()
}

func (b *BackendServer) Set(val bool) {
	b.isOnline.Store(val)
}

func NewBackendPool(settings []*BackendServerSettings) ([]*BackendServer, error) {
	backendList := make([]*BackendServer, len(settings))
	for i, b := range settings {
		if b.HelthUrl == "" {
			b.HelthUrl = "/health"
		}
		if b.Method == "" {
			b.Method = "GET"
		}

		urlB, err := url.Parse(b.BckndUrl)
		if err != nil {
			return nil, err
		}

		prx := httputil.NewSingleHostReverseProxy(urlB)

		backendList[i] = &BackendServer{
			BackendServerSettings: BackendServerSettings{
				BckndUrl: b.BckndUrl,
				Method:   b.Method,
				HelthUrl: b.BckndUrl + b.HelthUrl,
			},
			Prx: prx,
		}
	}

	return backendList, nil
}
