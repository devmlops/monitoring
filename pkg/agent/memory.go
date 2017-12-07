package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"log"
	"sort"
	"time"
)

type Memory struct {
	Time              time.Time       `json:"time"`
	MemoryTotalKB     uint64          `json:"memory_total_kb"`
	MemoryUsedKB      uint64          `json:"memory_used_kb"`
	MemoryUsedPercent float64         `json:"memory_used_percent"`
	MemoryByProcess   []ProcessMemory `json:"memory_by_process"`
}

type ProcessMemory struct {
	Pid               int32   `json:"pid"`
	Name              string  `json:"name"`
	MemoryUsedKB      uint64  `json:"memory_kb"`
	MemoryUsedPercent float32 `json:"memory_percent"`
}

func (m *Memory) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	m.GetMemoryUsageTotal()
	m.GetMemoryUsageByProcess()
}

func (m *Memory) GetMemoryUsageTotal() {
	m.Time = time.Now().UTC()
	stat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}
	m.MemoryTotalKB = stat.Total / 1024
	m.MemoryUsedKB = stat.Used / 1024
	m.MemoryUsedPercent = stat.UsedPercent
}

func (m *Memory) GetMemoryUsageByProcess() {
	reversed_freq := map[uint64][]ProcessMemory{}

	ps, _ := process.Processes()
	for _, proc := range ps {
		memPercent, err := proc.MemoryPercent()
		if err != nil {
			log.Fatal(err)
		}
		stat, err := proc.MemoryInfo()
		if err != nil {
			log.Fatal(err)
		}
		if stat.RSS > 0 {
			name, err := proc.Name()
			if err != nil {
				log.Fatal(err)
			}
			p := ProcessMemory{Name: name, Pid: proc.Pid, MemoryUsedPercent: memPercent, MemoryUsedKB: stat.RSS / 1024}
			reversed_freq[p.MemoryUsedKB] = append(reversed_freq[p.MemoryUsedKB], p)
		}
	}

	var numbers []int
	for val := range reversed_freq {
		numbers = append(numbers, int(val))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	for _, number := range numbers {
		for _, p := range reversed_freq[uint64(number)] {
			m.MemoryByProcess = append(m.MemoryByProcess, p)
		}
	}

	if Debug == true {
		ser, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(ser))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	res, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s", server.IP, server.port, "memory"),
		"application/json; charset=utf-8",
		b,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
