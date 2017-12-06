package main

import (
	"sync"
	"github.com/wwwthomson/monitoring/pkg/agent"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Analyse struct {
	store Store
	bot   *tgbotapi.BotAPI
	users []int64
}

type Store struct {
	netData []agent.Network
	memData []agent.Memory
	swData  []agent.Swap
	cpuData []agent.CPU

	mu sync.Mutex
}

func (s *Store) AddNetwork(n agent.Network) {
	s.mu.Lock()
	s.netData = append(s.netData, n)
	s.mu.Unlock()
	//go s.Analyse
	//SendAlert(s.bot, s.users, "test")
}

func (s *Store) AddMemory(n agent.Memory) {
	s.mu.Lock()
	s.memData = append(s.memData, n)
	s.mu.Unlock()
	//go s.Analyse
	//SendAlert(s.bot, s.users, "test")
}

func (s *Store) AddSwap(n agent.Swap) {
	s.mu.Lock()
	s.swData = append(s.swData, n)
	s.mu.Unlock()
	//go s.Analyse
	//SendAlert(s.bot, s.users, "test")
}

func (s *Store) AddCPU(n agent.CPU) {
	s.mu.Lock()
	s.cpuData = append(s.cpuData, n)
	s.mu.Unlock()
	//go s.Analyse
	//SendAlert(s.bot, s.users, "test")
}
