package main

import (
	"sync"
	"github.com/wwwthomson/monitoring/pkg/agent"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	//cpu := agent.CPU{}
	//go cpu.RunJob(&wg)

	network := agent.Network{}
	go network.RunJob(&wg)

	//swap := agent.Swap{}
	//go swap.RunJob(&wg)
    //
	//memory := agent.Memory{}
	//go memory.RunJob(&wg)
	wg.Wait()
}
