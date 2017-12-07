package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"log"
	"sort"
	"time"
)

type CPU struct {
	Time           time.Time    `json:"time"`
	CPUUsedPercent float64      `json:"cpu_used_percent"`
	CPUByProcess   []ProcessCPU `json:"cpu_by_process"`
	Server         Server       `json:"server"`
	Debug          bool         `json:"debug"`
	HostName       string       `json:"hostname"`
}

type ProcessCPU struct {
	Pid            int32   `json:"pid"`
	Name           string  `json:"name"`
	CPUUsedPercent float64 `json:"cpu_used_percent"`
}

func (c *CPU) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	c.GetCPUUsageTotal()
	c.GetCPUUsageByProcess()
}

func (c *CPU) GetCPUUsageTotal() {
	c.Time = time.Now().UTC()
	cpuUsage, err := cpu.Percent(time.Duration(1)*time.Second, false)
	if err != nil {
		log.Println(err)
	}
	c.CPUUsedPercent = cpuUsage[0]
}

func (c *CPU) GetCPUUsageByProcess() {
	c.CPUByProcess = nil
	reversed_freq := map[float64][]ProcessCPU{}

	ps, err := process.Processes()
	if err != nil {
		log.Println(err)
	}
	for _, pid := range ps {
		cpuPercent, err := pid.CPUPercent()
		if err != nil {
			log.Println(err)
		}
		if cpuPercent > 0 {
			name, err := pid.Name()
			if err != nil {
				log.Println(err)
			}
			if name != "hukumka-agent" {
				p := ProcessCPU{Name: name, Pid: pid.Pid, CPUUsedPercent: cpuPercent}
				reversed_freq[p.CPUUsedPercent] = append(reversed_freq[p.CPUUsedPercent], p)
			}
		}

	}

	var numbers []float64
	for val := range reversed_freq {
		numbers = append(numbers, val)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(numbers)))
	if len(numbers) > 5 {
		numbers = numbers[:5]
	}
	for _, number := range numbers {
		for _, p := range reversed_freq[number] {
			c.CPUByProcess = append(c.CPUByProcess, p)
		}
	}

	if c.Debug == true {
		ser, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(ser))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(c)
	res, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s", c.Server.IP, c.Server.Port, "cpu"),
		"application/json; charset=utf-8",
		b,
	)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
}
