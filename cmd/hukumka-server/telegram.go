package main

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/wwwthomson/monitoring/pkg/agent"
)

func RunTelegramBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("telegram bot: Authorized on account %s", bot.Self.UserName)
	return bot
}

func SendAlert(bot *tgbotapi.BotAPI, users []int64, message string) {
	for _, id := range users {
		ApiMessage := tgbotapi.NewMessage(id, message)
		ApiMessage.ParseMode = "markdown"
		bot.Send(ApiMessage)
	}
}

type FormMessageCPU struct {
	typeMessage string
	average     float64
	max         uint64
	real        float64
	message     string
	processes   []agent.ProcessCPU
	hostname    string
}

func (m *FormMessageCPU) SendAlertFromFormCPU(bot *tgbotapi.BotAPI, users []int64) {
	var message string
	message = fmt.Sprintf("**%s**: %s\n", m.typeMessage, m.hostname)
	message += fmt.Sprintf("%v\n", m.message)
	message += fmt.Sprintf("Среднее: %.2f\n", m.average)
	message += fmt.Sprintf("Максимальное: %v\n", m.max)
	message += fmt.Sprintf("Реальное: %.2f\n\n", m.real)
	if len(m.processes) != 0 {
		//<<<<<<< HEAD
		//		message += fmt.Sprintln("Top:")
		//		for i, process := range m.processes {
		//			k := i + 1
		//			message += fmt.Sprintln("%s: %s %s %s\n", k, process.Name, process.Pid, process.CPUUsedPercent)
		//=======
		message += fmt.Sprintf("Top processes:\n")
		for i, proc := range m.processes {
			k := i + 1
			message += fmt.Sprintf("%v: %s `%v` %.2f\n", k, proc.Name, proc.Pid, proc.CPUUsedPercent)
		}
	}
	SendAlert(bot, users, message)
}

type FormMessageNet struct {
	typeMessage string
	average     uint64
	max         uint64
	real        uint64
	message     string
	connections []agent.Connection
	hostname    string
}

type FormMessageMem struct {
	typeMessage   string
	average       uint64
	max           uint64
	real          uint64
	message       string
	processMemory []agent.ProcessMemory
	hostname      string
}

func (m *FormMessageNet) SendAlertFromFormNet(bot *tgbotapi.BotAPI, users []int64) {
	var message string
	message = fmt.Sprintf("**%s**: %s\n", m.typeMessage, m.hostname)
	message += fmt.Sprintf("%v\n", m.message)
	message += fmt.Sprintf("Среднее: %v\n", m.average)
	message += fmt.Sprintf("Максимальное: %v\n", m.max)
	message += fmt.Sprintf("Реальное: %v\n\n", m.real)
	if len(m.connections) != 0 {
		message += fmt.Sprintf("Top connections:\n")
		for i, connection := range m.connections {
			k := i + 1
			message += fmt.Sprintf("%v: `%s` %v\n", k, connection.IPAddress, connection.Number)
		}
	}
	//fmt.Println(">>> HERE\n\n")
	//fmt.Printf(message)
	SendAlert(bot, users, message)
}

func (m *FormMessageMem) SendAlertFromFormMem(bot *tgbotapi.BotAPI, users []int64) {
	var message string
	message = fmt.Sprintf("**%s**: %s\n", m.typeMessage, m.hostname)
	message += fmt.Sprintf("%v\n", m.message)
	message += fmt.Sprintf("Среднее: %v\n", m.average)
	message += fmt.Sprintf("Максимальное: %v\n", m.max)
	message += fmt.Sprintf("Реальное: %v\n\n", m.real)
	if len(m.processMemory) != 0 {
		message += fmt.Sprintf("Top connections:\n")
		for i, procMem := range m.processMemory {
			k := i + 1
			message += fmt.Sprintf("%v: `%s` %v\n", k, procMem.Name, procMem.MemoryUsedKB)
		}
	}
	SendAlert(bot, users, message)
}
