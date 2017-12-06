package main

func main() {
	config := OpenConfig("config.json")
	bot := RunTegelegramBot(config.TelegramBot.Token)
	//route := HttpServer(&Store{users: config.TelegramBot.Users, bot: bot})
	SendAlert(bot, config.TelegramBot.Users, "It's Ok" )
	monitor := Monitor{bot: bot, config: config}
	route := HttpServer(&monitor)
	route.Run() // listen and serve on 0.0.0.0:8080
}
