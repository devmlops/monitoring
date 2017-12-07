package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("telegram bot: Authorized on account %s", bot.Self.UserName)
	//return bot
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//updates, err := bot.GetUpdatesChan(u)
	//
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//	log.Printf(">>>>>> %s", update.Message.Chat.ID)
	//
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	bot.Send(msg)
	//}
	return bot
}

func SendAlert(bot *tgbotapi.BotAPI, users []int64, message string) {
	for _, id := range users {
		ApiMessage := tgbotapi.NewMessage(id, message)
		bot.Send(ApiMessage)
	}
}
