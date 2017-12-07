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

type Swap struct {
	Time            time.Time     `json:"time"`
	SwapTotalKB     uint64        `json:"swap_total_kb"`
	SwapUsedKB      uint64        `json:"swap_used_kb"`
	SwapUsedPercent float64       `json:"swap_used_percent"`
	SwapByProcess   []ProcessSwap `json:"swap_by_process"`
	Server          Server        `json:"server"`
	Debug           bool          `json:"debug"`
	HostName        string        `json:"hostname"`
}

type ProcessSwap struct {
	Pid             int32   `json:"pid"`
	Name            string  `json:"name"`
	SwapUsedKB      uint64  `json:"swap_kb"`
	SwapUsedPercent float64 `json:"swap_percent"`
}

func (s *Swap) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	s.GetSwapUsageTotal()
	s.GetSwapUsageByProcess()
}

func (s *Swap) GetSwapUsageTotal() {
	s.Time = time.Now().UTC()
	stat, err := mem.SwapMemory()
	if err != nil {
		log.Println(err)
	}
	s.SwapTotalKB = stat.Total / 1024
	s.SwapUsedKB = stat.Used / 1024
	s.SwapUsedPercent = stat.UsedPercent
}

func (s *Swap) GetSwapUsageByProcess() {
	s.SwapByProcess = nil
	reversed_freq := map[uint64][]ProcessSwap{}

	ps, err := process.Processes()
	if err != nil {
		log.Println(err)
	}
	for _, proc := range ps {
		stat, err := proc.MemoryInfo()
		if err != nil {
			log.Println(err)
		}
		if stat.Swap > 0 {
			used := stat.Swap / 1024
			name, err := proc.Name()
			if err != nil {
				log.Println(err)
			}
			swapPercent := float64(used) / float64(s.SwapTotalKB) * 100
			p := ProcessSwap{Name: name, Pid: proc.Pid, SwapUsedPercent: swapPercent, SwapUsedKB: used}
			reversed_freq[p.SwapUsedKB] = append(reversed_freq[p.SwapUsedKB], p)
		}
	}

	var numbers []int
	for val := range reversed_freq {
		numbers = append(numbers, int(val))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	if len(numbers) > 5 {
		numbers = numbers[:5]
	}
	for _, number := range numbers {
		for _, p := range reversed_freq[uint64(number)] {
			s.SwapByProcess = append(s.SwapByProcess, p)
		}
	}

	if s.Debug == true {
		ser, err := json.Marshal(s)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(ser))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(s)
	res, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s", s.Server.IP, s.Server.Port, "swap"),
		"application/json; charset=utf-8",
		b,
	)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
}
