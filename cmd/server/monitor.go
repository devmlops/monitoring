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
	netData  []agent.Network
	memData  []agent.Memory
	swData   []agent.Swap
	cpuData  []agent.CPU
	diskData []agent.Disk

	avarage Avarage

	Warning Status
	Danger  Status

	mu sync.Mutex
}

type Status struct {
	netStatus  bool
	memStatus  bool
	swStatus   bool
	cpuStatus  bool
	diskStatus bool

	netCounter  int
	memCounter  int
	swCounter   int
	cpuCounter  int
	diskCounter int
}

type Avarage struct {
	netAverage  float64
	memAverage  float64
	swAverage   float64
	cpuAverage  float64
	diskAverage float64
}

func (m *Monitor) AddNetwork(n agent.Network) {
	m.store.mu.Lock()
	m.store.netData = append(m.store.netData, n)
	m.store.mu.Unlock()
	//m.AnalyseNetwork(n)
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

func (m *Monitor) AddDisk(n agent.Disk) {
	m.store.mu.Lock()
	m.store.diskData = append(m.store.diskData, n)
	m.store.mu.Unlock()
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddCPU(n agent.CPU) {
	m.store.mu.Lock()
	m.store.cpuData = append(m.store.cpuData, n)
	m.store.avarage.cpuAverage += n.CPUUsedPercent
	m.store.mu.Unlock()
	m.AnalyseCPU(n)
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AnalyseCPU(n agent.CPU) {
	m.store.mu.Lock()
	result := m.store.avarage.cpuAverage / float64(len(m.store.netData))
	m.store.mu.Unlock()
	// +20%
	if n.CPUUsedPercent > result*1.2 {
		if n.CPUUsedPercent >= float64(m.config.CPU.MaxLimit) && m.store.Danger.cpuStatus == false {
			// check danger
			if m.store.Danger.cpuCounter == 3 && m.store.Danger.cpuStatus != true {
				m.store.Warning.cpuStatus = true
				m.store.Danger.cpuStatus = true
				m.store.Warning.cpuCounter = 3
				//fmt.Println("Allert Danger")
				SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.cpuCounter += 1
			}
		} else {
			// check warning
			if m.store.Warning.cpuCounter == 3 && m.store.Warning.cpuStatus != true {
				m.store.Warning.cpuStatus = true
				//fmt.Println("Allert Warning")
				SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
				//SendAlert(m.bot, m.config.TelegramBot.Users, fmt.Sprintf("%s \n среднее: %v \n реальные: %v", "Network: превышение свыше 20%", int(result), int(n.Connections)))
			} else {
				m.store.Warning.cpuCounter += 1
			}
		}
	} else {
		// normalization
		if m.store.Danger.cpuCounter != 0 {
			m.store.Danger.cpuCounter -= 1
		}
		if m.store.Warning.cpuCounter != 0 {
			m.store.Warning.cpuCounter -= 1
		}
		if m.store.Warning.cpuStatus == true && m.store.Danger.cpuCounter == 0 && m.store.Warning.cpuCounter == 0 {
			//fmt.Println("All good")
			m.store.Danger.cpuStatus = false
			m.store.Warning.cpuStatus = false
			SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
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
