package agent

import (
	"sort"
	"fmt"
	"time"
	"sync"
	//"bytes"
	//"net/http"
	"encoding/json"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type Swap struct {
	Time            time.Time     `json:"time"`
	SwapTotalKB     uint64           `json:"swap_total_kb"`
	SwapUsedKB      uint64           `json:"swap_used_kb"`
	SwapUsedPercent float64       `json:"swap_used_percent"`
	SwapByProcess   []ProcessSwap `json:"swap_by_process"`
}

type ProcessSwap struct {
	Pid         int32     `json:"pid"`
	Name        string  `json:"name"`
	SwapUsedKB      uint64     `json:"swap_kb"`
	SwapUsedPercent float64 `json:"swap_percent"`
}

func (s *Swap) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	s.GetSwapUsageTotal()
	s.GetSwapUsageByProcess()
}

func (s *Swap) GetSwapUsageTotal() {
	stat, _ := mem.VirtualMemory()
	s.SwapTotalKB = stat.Total / 1024
	s.SwapUsedKB = stat.Used / 1024
	s.SwapUsedPercent = stat.UsedPercent
	//fmt.Println(s)
}

func (s *Swap) GetSwapUsageByProcess() {
	s.Time = time.Now().UTC()
	
	reversed_freq := map[uint64][]ProcessSwap{}
	
	ps, _ := process.Processes()
	for _, proc := range ps {
		stat, _ := proc.MemoryInfo()
		fmt.Println(stat.Swap)
		used := stat.Swap/1024
		name, _ := proc.Name()
		swapPercent := float64(used) /  float64(s.SwapTotalKB) * 100
		p := ProcessSwap{Name: name, Pid: proc.Pid, SwapUsedPercent: swapPercent, SwapUsedKB: used}
		reversed_freq[p.SwapUsedKB] = append(reversed_freq[p.SwapUsedKB], p)
	}
	
	var numbers []int
	for val := range reversed_freq {
		numbers = append(numbers, int(val))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	for _, number := range numbers {
		for _, p := range reversed_freq[uint64(number)] {
			s.SwapByProcess = append(s.SwapByProcess, p)
		}
	}
	
	ser, _ := json.Marshal(s)
	fmt.Println(string(ser))
	
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(s)
	res, _ := http.Post("http://192.168.88.141:8080/memory", "application/json; charset=utf-8", b)
	fmt.Println(res)
}


