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
	Current uint64
	Mu      sync.RWMutex
	Pool    []*BackendServer
}
