package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
	"sync"
	"time"
	"fmt"
)

func main() {
	config := OpenConfig("config.json")
	debug := config.Debug
	server := config.Server
	
	cpu := agent.CPU{Server: server, Debug: debug}
	memory := agent.Memory{Server: server, Debug: debug}
	swap := agent.Swap{Server: server, Debug: debug}
	network := agent.Network{Server: server, Debug: debug}
	disk := agent.Disk{Server: server, Debug: debug}
	var wg sync.WaitGroup

	if debug == true {
		wg.Add(5)
		p := agent.Params{UseWg: true, Wg: &wg}

		go cpu.RunJob(&p)
		go memory.RunJob(&p)
		go swap.RunJob(&p)
		go network.RunJob(&p)
		go disk.RunJob(&p)

		wg.Wait()
	} else {
		for {
			p := agent.Params{UseWg: false, Wg: &wg}

			go cpu.RunJob(&p)
			go memory.RunJob(&p)
			go swap.RunJob(&p)
			go network.RunJob(&p)
			go disk.RunJob(&p)
			time.Sleep(1 * time.Second)
		}
	}
}
