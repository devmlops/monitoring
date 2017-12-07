package main

func main() {
	config := OpenConfig("config.json")
	bot := RunTelegramBot(config.TelegramBot.Token)
	//SendAlert(bot, config.TelegramBot.Users, "It's Ok" )
	monitor := Monitor{bot: bot, config: config}
	route := HttpServer(&monitor)
	route.Run(config.Server.Address) // listen and serve on 0.0.0.0:8080
}
