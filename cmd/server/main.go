package main

func main() {
	bot := RunTegelegramBot("test token")
	SendAlert("test", bot)
}
