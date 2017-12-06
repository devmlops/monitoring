package main

import (
	"sync"
	"github.com/wwwthomson/monitoring/pkg/agent"
)

type Store struct {
	netData []agent.Network
	memData []agent.Memory
	swData  []agent.Swap
	cpuData []agent.CPU

	mu sync.Mutex
}

func (s *Store) AddNetwork(n agent.Network)  {
	s.mu.Lock()
	s.netData = append(s.netData, n)
	s.mu.Unlock()
	//go s.Analyse
}

func (s *Store) AddMemory(n agent.Memory)  {
	s.mu.Lock()
	s.memData = append(s.memData, n)
	s.mu.Unlock()
	//go s.Analyse
}