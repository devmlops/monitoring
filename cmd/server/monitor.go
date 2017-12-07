package main

import (
	"sync"
	"github.com/wwwthomson/monitoring/pkg/agent"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
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
	m.AnalyseNetwork(n)
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddMemory(n agent.Memory) {
	m.store.mu.Lock()
	m.store.memData = append(m.store.memData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddSwap(n agent.Swap) {
	m.store.mu.Lock()
	m.store.swData = append(m.store.swData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddCPU(n agent.CPU) {
	m.store.mu.Lock()
	m.store.cpuData = append(m.store.cpuData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AnalyseNetwork(n agent.Network) {
	fmt.Println(n.Time)
	fmt.Println(n.Connections)
	fmt.Println(n.ConnectionsByIP)
	var total uint64
	for _, value := range m.store.netData {
		total += value.Connections
	}
	result := float64(total) / float64(len(m.store.netData))
	// +20%
	if float64(n.Connections) > result * 1.2 {
		//fmt.Sprintf("%s %s/%s", "Network: превышение на 20%", int(result), int(n.Connections))
		SendAlert(m.bot, m.config.TelegramBot.Users, fmt.Sprintf("%s \n среднее: %v \n реальные: %v", "Network: превышение свыше 20%", int(result), int(n.Connections)))
	}
	if n.Connections >= m.config.Network.MaxLimit {
		SendAlert(m.bot, m.config.TelegramBot.Users, "(_*_)")
	}
}

func (m *Monitor) AnalyseMemory(n agent.Memory) {
	//fmt.Println(n.Time)
	//mean := GetMean()
	fmt.Println(m.store.netData[len(m.store.netData)-1].Connections)
}
//
//func GetMean(massive []uint64) float64 {
//	var total uint64
//	for _, value := range massive {
//		total += value
//	}
//	result := float64(total) / float64(len(massive))
//	return result
//}