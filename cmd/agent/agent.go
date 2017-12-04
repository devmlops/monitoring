package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
)

func main() {
	//network := agent.Network{}
	//network.RunJob()
    //
	//swap := agent.Swap{}
	//swap.RunJob()

	memory := agent.Memory{}
	memory.RunJob()
}
