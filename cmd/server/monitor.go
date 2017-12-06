package main

import (
	"sync"
	"github.com/wwwthomson/monitoring/pkg/agent"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Monitor struct {
	store  Store
	bot    *tgbotapi.BotAPI
	config Config
}

type Store struct {
	netData []agent.Network
	memData []agent.Memory
	swData  []agent.Swap
	cpuData []agent.CPU

	mu sync.Mutex
}

func (m *Monitor) AddNetwork(n agent.Network) {
	m.store.mu.Lock()
	m.store.netData = append(m.store.netData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddMemory(n agent.Memory) {
	m.store.mu.Lock()
	m.store.memData = append(m.store.memData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddSwap(n agent.Swap) {
	m.store.mu.Lock()
	m.store.swData = append(m.store.swData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddCPU(n agent.CPU) {
	m.store.mu.Lock()
	m.store.cpuData = append(m.store.cpuData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}
