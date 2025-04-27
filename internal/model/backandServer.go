package model

import (
	"net/http/httputil"
	"sync"
)

type BackendServer struct {
	IsOnline bool
	Url    string
	Prx    *httputil.ReverseProxy
}

type BackendPool struct {
	Mu   sync.Mutex
	Pool []*BackendServer
}
