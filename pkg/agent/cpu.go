package agent

import (
	"log"
	"fmt"
	"time"
	"sync"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type CPU struct {
	Time              time.Time       `json:"time"`
	CPUUsedPercent    float64         `json:"cpu_used_percent"`
	CPUByProcess      []ProcessCPU    `json:"cpu_by_process"`
}

type ProcessCPU struct {
	PID           int     `json:"pid"`
	Name          string  `json:"name"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
}

func (c *CPU) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	c.GetCPUUsageTotal()
	//c.GetCPUUsageByProcess()
}

func (c *CPU) GetCPUUsageTotal() {
	c.Time = time.Now().UTC()
	cpuUsage, err := cpu.Percent(time.Duration(1)*time.Second, false)
	if err != nil {
		log.Fatal(err)
	}
	c.CPUUsedPercent = cpuUsage[0]
	ser, err := json.Marshal(c)
	fmt.Println(string(ser))
}

func (c *CPU) GetCPUUsageByProcess() {
	ps, _ := process.Processes()
	fmt.Println(ps[0])
	for _, pid := range ps {
		cpuPercent, _ := pid.CPUPercent()
		name, _ := pid.Name()
		fmt.Println(name)
		fmt.Println(cpuPercent)
	}
}
