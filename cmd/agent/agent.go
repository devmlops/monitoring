package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
)

func main() {
	network := agent.Network{}
	network.GetActiveConnections()

	swap := agent.Swap{}
	swap.GetSwapUsageTotal()
	swap.GetSwapUsageByProcess()

	memory := agent.Memory{}
	memory.GetMemoryUsageTotal()
	memory.GetMemoryUsageByProcess()
}
