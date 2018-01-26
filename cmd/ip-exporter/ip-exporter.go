package main

import (
	ip_exporter "github.com/wwwthomson/monitoring/pkg/ip-exporter"
	"sync"
)

func main() {
	
	var wg sync.WaitGroup
    
	conn := ip_exporter.Connections{}
    
    wg.Add(1)
	
	p := ip_exporter.Params{UseWg: true, Wg: &wg}
	
	go conn.RunJob(&p)
    
    wg.Wait()
}
