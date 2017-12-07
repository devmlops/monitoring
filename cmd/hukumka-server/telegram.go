package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot(token string) *tgbotapi.BotAPI{
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
