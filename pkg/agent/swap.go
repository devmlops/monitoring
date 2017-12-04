package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"sync"
)

type Swap struct {
	Time            time.Time     `json:"time"`
	SwapTotalKB     int           `json:"swap_total_kb"`
	SwapUsedKB      int           `json:"swap_used_kb"`
	SwapUsedPercent float32       `json:"swap_used_percent"`
	SwapByProcess   []ProcessSwap `json:"swap_by_process"`
}

type ProcessSwap struct {
	PID         int     `json:"pid"`
	Name        string  `json:"name"`
	SwapKB      int     `json:"swap_kb"`
	SwapPercent float32 `json:"swap_percent"`
}

func (s *Swap) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	s.GetSwapUsageTotal()
	s.GetSwapUsageByProcess()
}

func (s *Swap) GetSwapUsageByProcess() {
	s.Time = time.Now().UTC()
	swapCmd := `for file in /proc/*/status ; do awk '/^Pid|VmSwap|Name/{printf $2 " "}END{ print ""}' $file; done | sort -k 3 -n -r`
	cmd := exec.Command(
		"/bin/bash",
		"-c",
		swapCmd,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		perProcess := strings.Fields(line)
		if len(perProcess) > 2 {
			PID, err := strconv.Atoi(perProcess[1])
			if err != nil {
				log.Fatal(err)
			}
			swap, err := strconv.Atoi(perProcess[2])
			if err != nil {
				log.Fatal(err)
			}
			if swap > 0 {
				p := ProcessSwap{Name: perProcess[0], PID: PID, SwapKB: swap}
				p.SwapPercent = float32(p.SwapKB) / float32(s.SwapTotalKB) * 100.0
				s.SwapByProcess = append(s.SwapByProcess, p)
			}
		}
	}
	ser, err := json.Marshal(s)
	fmt.Println(string(ser))
}

func (s *Swap) GetSwapUsageTotal() {
	swapCmd := `cat /proc/swaps | tail -n1`
	cmd := exec.Command(
		"/bin/bash",
		"-c",
		swapCmd,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	swap := strings.Fields(out.String())
	s.SwapTotalKB, err = strconv.Atoi(swap[2])
	if err != nil {
		log.Fatal(err)
	}
	s.SwapUsedKB, err = strconv.Atoi(swap[3])
	if err != nil {
		log.Fatal(err)
	}
	s.SwapUsedPercent = float32(s.SwapUsedKB) / float32(s.SwapTotalKB) * 100.0
}
