package main

import (
	"github.com/wwwthomson/monitoring/pkg/agent"
	"sync"
	//"fmt"
	"time"
)

func main() {
	//cpu := agent.CPU{}
	//memory := agent.Memory{}
	//swap := agent.Swap{}
	//network := agent.Network{}
	disk := agent.Disk{}
	
	if agent.Debug == true {
		var wg sync.WaitGroup
		
		wg.Add(1)
		
		p := agent.Params{UseWg: true, Wg: &wg}
		
		//go cpu.RunJob(&wg)
		//go memory.RunJob(&wg)
		//go swap.RunJob(&wg)
		//go network.RunJob(&wg)
		
		go disk.RunJob(&p)

		wg.Wait()
	} else {
		for {
			//go cpu.RunJob(&wg)
			//go memory.RunJob(&wg)
			//go swap.RunJob(&wg)
			//go network.RunJob(&wg)
			//go disk.RunJob(&wg)
			time.Sleep(5 * time.Second)
		}
	}
}
