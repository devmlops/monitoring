package agent

import (
	"log"
	"fmt"
	"time"
	"sync"
	"bytes"
	"net/http"
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
	Pid           int32     `json:"pid"`
	Name          string  `json:"name"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
}

func (c *CPU) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	c.GetCPUUsageTotal()
	c.GetCPUUsageByProcess()
}

func (c *CPU) GetCPUUsageTotal() {
	c.Time = time.Now().UTC()
	cpuUsage, err := cpu.Percent(time.Duration(1)*time.Second, false)
	if err != nil {
		log.Fatal(err)
	}
	c.CPUUsedPercent = cpuUsage[0]
}

func (c *CPU) GetCPUUsageByProcess() {
	ps, _ := process.Processes()
	for _, pid := range ps {
		cpuPercent, _ := pid.CPUPercent()
		name, _ := pid.Name()
		p := ProcessCPU{Name: name, Pid: pid.Pid, CPUUsedPercent: cpuPercent}
		c.CPUByProcess = append(c.CPUByProcess, p)
	}
	ser, _ := json.Marshal(c)
	fmt.Println(string(ser))
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(c)
	res, _ := http.Post("http://192.168.88.141:8080/cpu", "application/json; charset=utf-8", b)
	fmt.Println(res)
}
