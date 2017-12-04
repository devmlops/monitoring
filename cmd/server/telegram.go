package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTegelegramBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("telegram bot: Authorized on account %s", bot.Self.UserName)
	return bot
}


func SendAlert(message string, bot *tgbotapi.BotAPI) {
	Users := []int64{ 282049937 }
	for _, id := range Users {
		ApiMessage := tgbotapi.NewMessage(id, message)
		bot.Send(ApiMessage)
	}
}