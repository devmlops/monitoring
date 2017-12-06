package main

func main() {
	//bot := RunTegelegramBot("test token")
	//SendAlert("test", bot)
	route := HttpServer(&Store{})
	route.Run() // listen and serve on 0.0.0.0:8080
}