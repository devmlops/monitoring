package main

import (
	//"fmt"
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
	netData  []agent.Network
	memData  []agent.Memory
	swapData []agent.Swap
	cpuData  []agent.CPU
	diskData []agent.Disk

	average average

	Warning Status
	Danger  Status

	mu sync.Mutex
}

type Status struct {
	netStatus  bool
	memStatus  bool
	swapStatus bool
	cpuStatus  bool
	diskStatus bool

	netCounter  int
	memCounter  int
	swapCounter int
	cpuCounter  int
	diskCounter int
}

type average struct {
	netAverage  uint64
	memAverage  uint64
	swapAverage uint64
	cpuAverage  float64
	diskAverage uint64
}

func (m *Monitor) AddNetwork(n agent.Network) {
	//fmt.Println(">>> HERE 3")
	m.store.mu.Lock()
	m.store.netData = append(m.store.netData, n)
	m.store.average.netAverage += n.Connections
	m.store.mu.Unlock()
	m.AnalyseNetwork(n)
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddMemory(n agent.Memory) {
	m.store.mu.Lock()
	m.store.memData = append(m.store.memData, n)
	m.store.average.memAverage += n.MemoryUsedKB
	m.store.mu.Unlock()
	m.AnalyseMemory(n)
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddSwap(n agent.Swap) {
	m.store.mu.Lock()
	m.store.swapData = append(m.store.swapData, n)
	m.store.average.swapAverage += n.SwapUsedKB
	m.store.mu.Unlock()
	m.AnalyseSwap(n)
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddDisk(n agent.Disk) {
	m.store.mu.Lock()
	m.store.diskData = append(m.store.diskData, n)
	m.store.average.diskAverage += n.DiskUsedKB
	m.store.mu.Unlock()
	m.AnalyseDisk(n)
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AddCPU(n agent.CPU) {
	m.store.mu.Lock()
	m.store.cpuData = append(m.store.cpuData, n)
	m.store.average.cpuAverage += n.CPUUsedPercent
	m.store.mu.Unlock()
	m.AnalyseCPU(n)
	//go m.Monitor
	//SendAlert(m.bot, m.config.TelegramBot.Users, "test")
}

func (m *Monitor) AnalyseNetwork(n agent.Network) {
	m.store.mu.Lock()
	//fmt.Println(">>> HERE 6")
	result := float64(m.store.average.netAverage) / float64(len(m.store.netData))
	// +20%
	if float64(n.Connections) > float64(result*1.2) {
	//if true {
		if n.Connections >= m.config.Network.MaxLimit && m.store.Danger.netStatus == false {
		//if true {
			// check danger
			if m.store.Danger.netCounter == 3 && m.store.Danger.netStatus != true {
			//if true {
				m.store.Warning.netStatus = true
				m.store.Danger.netStatus = true
				m.store.Warning.netCounter = 3
				//fmt.Println("Allert Danger")
				fm := FormMessageNet{
					typeMessage: "DANGER",
					average: m.store.average.netAverage,
					max: m.config.Network.MaxLimit,
					real: m.store.netData[0].Connections,
					connections: m.store.netData[0].ConnectionsByIP,
					hostname: m.store.netData[0].Hostname,
				}
				//fmt.Println(fm)
				go fm.SendAlertFromFormNet(m.bot, m.config.TelegramBot.Users)
				//go SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.netCounter += 1
			}
		} else {
			if m.store.Warning.netCounter == 3 && m.store.Warning.netStatus != true {
				m.store.Warning.netStatus = true
				fm := FormMessageNet{
					typeMessage: "WARNING",
					average: m.store.average.netAverage,
					max: m.config.Network.MaxLimit,
					real: m.store.netData[0].Connections,
					connections: m.store.netData[0].ConnectionsByIP,
					hostname: m.store.netData[0].Hostname,
				}
				//fmt.Println(fm)
				go fm.SendAlertFromFormNet(m.bot, m.config.TelegramBot.Users)
				//go SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
			} else {
				m.store.Warning.netCounter += 1
			}
		}
	} else {
		// normalization
		if m.store.Danger.netCounter != 0 {
			m.store.Danger.netCounter -= 1
		}
		if m.store.Warning.netCounter != 0 {
			m.store.Warning.netCounter -= 1
		}
		if m.store.Warning.netStatus == true && m.store.Danger.netCounter == 0 && m.store.Warning.netCounter == 0 {
			//fmt.Println("All good")
			m.store.Danger.netStatus = false
			m.store.Warning.netStatus = false
			go SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
	}
	m.store.mu.Unlock()
}

func (m *Monitor) AnalyseMemory(n agent.Memory) {
	m.store.mu.Lock()
	result := float64(m.store.average.memAverage) / float64(len(m.store.memData))
	// +20%
	if float64(n.MemoryUsedKB) > result*1.2 {
		if n.MemoryUsedKB >= m.config.Memory.MaxLimit && m.store.Danger.memStatus == false {
			// check danger
			if m.store.Danger.memCounter == 3 && m.store.Danger.memStatus != true {
				m.store.Warning.memStatus = true
				m.store.Danger.memStatus = true
				m.store.Warning.memCounter = 3
				//fmt.Println("Allert Danger")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.memCounter += 1
			}
		} else {
			// check warning
			if m.store.Warning.memCounter == 3 && m.store.Warning.memStatus != true {
				m.store.Warning.memStatus = true
				//fmt.Println("Allert Warning")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
				//SendAlert(m.bot, m.config.TelegramBot.Users, fmt.Sprintf("%s \n среднее: %v \n реальные: %v", "Network: превышение свыше 20%", int(result), int(n.Connections)))
			} else {
				m.store.Warning.memCounter += 1
			}
		}
	} else {
		// normalization
		if m.store.Danger.memCounter != 0 {
			m.store.Danger.memCounter -= 1
		}
		if m.store.Warning.memCounter != 0 {
			m.store.Warning.memCounter -= 1
		}
		if m.store.Warning.memStatus == true && m.store.Danger.memCounter == 0 && m.store.Warning.memCounter == 0 {
			//fmt.Println("All good")
			m.store.Danger.memStatus = false
			m.store.Warning.memStatus = false
			go SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
	}
	m.store.mu.Unlock()
}

func (m *Monitor) AnalyseSwap(n agent.Swap) {
	m.store.mu.Lock()
	result := float64(m.store.average.swapAverage) / float64(len(m.store.swapData))
	// +20%
	if float64(n.SwapUsedKB) > result*1.2 {
		if n.SwapUsedKB >= m.config.Swap.MaxLimit && m.store.Danger.swapStatus == false {
			// check danger
			if m.store.Danger.swapCounter == 3 && m.store.Danger.swapStatus != true {
				m.store.Warning.swapStatus = true
				m.store.Danger.swapStatus = true
				m.store.Warning.swapCounter = 3
				//fmt.Println("Allert Danger")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.swapCounter += 1
			}
		} else {
			// check warning
			if m.store.Warning.swapCounter == 3 && m.store.Warning.swapStatus != true {
				m.store.Warning.swapStatus = true
				//fmt.Println("Allert Warning")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
				//SendAlert(m.bot, m.config.TelegramBot.Users, fmt.Sprintf("%s \n среднее: %v \n реальные: %v", "Network: превышение свыше 20%", int(result), int(n.Connections)))
			} else {
				m.store.Warning.swapCounter += 1
			}
		}
	} else {
		// normalization
		if m.store.Danger.swapCounter != 0 {
			m.store.Danger.swapCounter -= 1
		}
		if m.store.Warning.swapCounter != 0 {
			m.store.Warning.swapCounter -= 1
		}
		if m.store.Warning.swapStatus == true && m.store.Danger.swapCounter == 0 && m.store.Warning.swapCounter == 0 {
			//fmt.Println("All good")
			m.store.Danger.swapStatus = false
			m.store.Warning.swapStatus = false
			go SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
	}
	m.store.mu.Unlock()
}

func (m *Monitor) AnalyseDisk(n agent.Disk) {
	m.store.mu.Lock()
	result := float64(m.store.average.diskAverage) / float64(len(m.store.diskData))
	// +20%
	if float64(n.DiskUsedKB) > result*1.2 {
		if n.DiskUsedKB >= m.config.Disk.MaxLimit && m.store.Danger.diskStatus == false {
			// check danger
			if m.store.Danger.diskCounter == 3 && m.store.Danger.diskStatus != true {
				m.store.Warning.diskStatus = true
				m.store.Danger.diskStatus = true
				m.store.Warning.diskCounter = 3
				//fmt.Println("Allert Danger")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.diskCounter += 1
			}
		} else {
			// check warning
			if m.store.Warning.diskCounter == 3 && m.store.Warning.diskStatus != true {
				m.store.Warning.diskStatus = true
				//fmt.Println("Allert Warning")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
				//SendAlert(m.bot, m.config.TelegramBot.Users, fmt.Sprintf("%s \n среднее: %v \n реальные: %v", "Network: превышение свыше 20%", int(result), int(n.Connections)))
			} else {
				m.store.Warning.diskCounter += 1
			}
		}
	} else {
		// normalization
		if m.store.Danger.diskCounter != 0 {
			m.store.Danger.diskCounter -= 1
		}
		if m.store.Warning.diskCounter != 0 {
			m.store.Warning.diskCounter -= 1
		}
		if m.store.Warning.diskStatus == true && m.store.Danger.diskCounter == 0 && m.store.Warning.diskCounter == 0 {
			//fmt.Println("All good")
			m.store.Danger.diskStatus = false
			m.store.Warning.diskStatus = false
			go SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
	}
	m.store.mu.Unlock()
}

//func SendCPU {}

func (m *Monitor) AnalyseCPU(n agent.CPU) {
	m.store.mu.Lock()
	result := m.store.average.cpuAverage / float64(len(m.store.netData))
	// +20%
	if n.CPUUsedPercent > result*1.2 {
		if n.CPUUsedPercent >= float64(m.config.CPU.MaxLimit) && m.store.Danger.cpuStatus == false {
			// check danger
			if m.store.Danger.cpuCounter == 3 && m.store.Danger.cpuStatus != true {
				m.store.Warning.cpuStatus = true
				m.store.Danger.cpuStatus = true
				m.store.Warning.cpuCounter = 3
				//fmt.Println("Allert Danger")

				//message := FormMessage{
				//	typeMessage: "Danger",
				//	from:        "CPU",
				//	average:     result,
				//	max:         m.config.CPU.MaxLimit,
				//	real:        n.CPUUsedPercent,
				//	message:     "Достигнут максимальный лимит",
				//	processes:   func string{return "bla"}()
				//}

				go SendAlert(m.bot, m.config.TelegramBot.Users, "Danger")
			} else {
				m.store.Danger.cpuCounter += 1
			}
		} else {
			// check warning
			if m.store.Warning.cpuCounter == 3 && m.store.Warning.cpuStatus != true {
				m.store.Warning.cpuStatus = true
				//fmt.Println("Allert Warning")
				go SendAlert(m.bot, m.config.TelegramBot.Users, "Warning")
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
			go SendAlert(m.bot, m.config.TelegramBot.Users, "All good")
		}
	}
	m.store.mu.Unlock()
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
