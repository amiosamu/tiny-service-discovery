package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	proxy       *httputil.ReverseProxy
	containerID string
}

type ServiceRegistry struct {
	BackendsStore atomic.Value
}

func (s *ServiceRegistry) Init() {
	s.BackendsStore.Store([]Backend{})
}

func (s *ServiceRegistry) Add(containerID, addr string) {
	url, _ := url.Parse(addr)

	s.BackendsStore.Swap(append(s.BackendsStore.Load().([]Backend), Backend{
		proxy:       httputil.NewSingleHostReverseProxy(url),
		containerID: containerID,
	}))
}





func (s *ServiceRegistry) GetByContainerID(containerID string) (Backend, bool) {
	for _, b := range s.GetBackends() {
		if b.containerID == containerID {
			return b, true
		}
	}

	return Backend{}, false
}



func (s *ServiceRegistry) GetByIndex(index int) Backend {
	return s.GetBackends()[index]
}

func (s *ServiceRegistry) RemoveByContainerID(containerID string) {
	var backends []Backend
	for _, b := range s.GetBackends() {
		if b.containerID == containerID {
			continue
		}
		backends = append(backends, b)
	}

	s.BackendsStore.Store(backends)
}


func (s *ServiceRegistry) List() {
	backendList := s.GetBackends()
	for i := range backendList {
		fmt.Printf("Backend %d: %s\n", i, backendList[i].containerID)
	}
}

func (s *ServiceRegistry) RemoveAll() {
	s.BackendsStore.Store([]Backend{})
}

func (s *ServiceRegistry) Len() int {
	return len(s.GetBackends())
}

func (s *ServiceRegistry) GetBackends() []Backend {
	return s.BackendsStore.Load().([]Backend)
}
