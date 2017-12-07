package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
	"sync"
	"time"
	"os"
	"log"
)

func main() {
	config := OpenConfig("config.json")
	debug := config.Debug
	server := config.Server
	
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	cpu := agent.CPU{Server: server, Debug: debug, Hostname: hostname}
	memory := agent.Memory{Server: server, Debug: debug, Hostname: hostname}
	swap := agent.Swap{Server: server, Debug: debug, Hostname: hostname}
	network := agent.Network{Server: server, Debug: debug, Hostname: hostname}
	disk := agent.Disk{Server: server, Debug: debug, Hostname: hostname}
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
