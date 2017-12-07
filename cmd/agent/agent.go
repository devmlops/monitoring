package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
	"sync"
	"time"
)

func main() {
	cpu := agent.CPU{}
	memory := agent.Memory{}
	swap := agent.Swap{}
	network := agent.Network{}
	disk := agent.Disk{}
	var wg sync.WaitGroup

	if agent.Debug == true {
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
			time.Sleep(5 * time.Second)
		}
	}
}
