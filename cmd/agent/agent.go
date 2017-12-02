package main

import (
	"fmt"
	"github.com/wwwthomson/monitoring/pkg/agent"
)

func main() {
	network := agent.Network{}
	//network.Time = time.Now()
	network.GetActiveConnections()
	fmt.Println(network)
}

