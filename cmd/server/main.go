package main

func main() {
	config := OpenConfig("config.json")
	bot := RunTelegramBot(config.TelegramBot.Token)
	//RunTelegramBot(config.TelegramBot.Token)
	//route := HttpServer(&Store{users: config.TelegramBot.Users, bot: bot})
	//SendAlert(bot, config.TelegramBot.Users, "It's Ok" )
	monitor := Monitor{bot: bot, config: config}
	route := HttpServer(&monitor)
	route.Run() // listen and serve on 0.0.0.0:8080
	//x := []uint64{ 1, 2, 3, 4, 5 }
	//fmt.Println(GetMean(x))
}
