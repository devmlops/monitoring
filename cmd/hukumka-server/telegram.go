package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
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
		bot.Send(ApiMessage)
	}
}

type FormMessage struct {
	typeMessage string
	from        string
	avarage     float64
	max         uint64
	real        uint64
	message     string
	processes     []string
	server agent.Server
}

func (m *FormMessage) SendAlertFromForm(bot *tgbotapi.BotAPI, users []int64) {
	var message string
	message = fmt.Sprintln("**%s**:", m.typeMessage, m.server.IP)
	message += fmt.Sprintln("%s", m.from)
	message += fmt.Sprintln("%s", m.message)
	message += fmt.Sprintln("Среднее: %s", m.avarage)
	message += fmt.Sprintln("Максимальное: %s", m.max)
	message += fmt.Sprintln("Реальное: %s", m.real)
	if len(m.processes) != 0 {
		message += fmt.Sprintln("Top:")
		for process := range m.processes {
			message += fmt.Sprintln(process)
		}
	}
	SendAlert(bot, users, message)
}