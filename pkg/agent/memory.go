package agent

import (
	//"log"
	"sort"
	"fmt"
	"time"
	"sync"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type Memory struct {
	Time              time.Time       `json:"time"`
	MemoryTotalKB     uint64             `json:"memory_total_kb"`
	MemoryUsedKB      uint64             `json:"memory_used_kb"`
	MemoryUsedPercent float64         `json:"memory_used_percent"`
	MemoryByProcess   []ProcessMemory `json:"memory_by_process"`
}

type ProcessMemory struct {
	Pid           int32     `json:"pid"`
	Name          string  `json:"name"`
	MemoryUsedKB      uint64     `json:"memory_kb"`
	MemoryUsedPercent float32 `json:"memory_percent"`
}

func (m *Memory) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	m.GetMemoryUsageTotal()
	m.GetMemoryUsageByProcess()
}

func (m *Memory) GetMemoryUsageTotal() {
	m.Time = time.Now().UTC()
	stat, _ := mem.VirtualMemory()
	m.MemoryTotalKB = stat.Total / 1024
	m.MemoryUsedKB = stat.Used / 1024
	m.MemoryUsedPercent = stat.UsedPercent
}

func (m *Memory) GetMemoryUsageByProcess() {
	reversed_freq := map[uint64][]ProcessMemory{}
	
	ps, _ := process.Processes()
	for _, proc := range ps {
		memPercent, _ := proc.MemoryPercent()
		stat, _ := proc.MemoryInfo()
		if stat.RSS > 0 {
			name, _ := proc.Name()
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
	
	ser, _ := json.Marshal(m)
	fmt.Println(string(ser))
	
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	res, _ := http.Post("http://192.168.88.141:8080/memory", "application/json; charset=utf-8", b)
	fmt.Println(res)
}
