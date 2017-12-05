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

type Memory struct {
	Time              time.Time       `json:"time"`
	MemoryTotalKB     int             `json:"memory_total_kb"`
	MemoryUsedKB      int             `json:"memory_used_kb"`
	MemoryUsedPercent float32         `json:"memory_used_percent"`
	MemoryByProcess   []ProcessMemory `json:"memory_by_process"`
}

type ProcessMemory struct {
	PID           int     `json:"pid"`
	Name          string  `json:"name"`
	MemoryKB      int     `json:"memory_kb"`
	MemoryPercent float32 `json:"memory_percent"`
}

func (m *Memory) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	m.GetMemoryUsageTotal()
	m.GetMemoryUsageByProcess()
}

func (m *Memory) GetMemoryUsageByProcess() {
	m.Time = time.Now().UTC()
	memoryCmd := `for file in /proc/*/status ; do awk '/^Pid|VmRSS|Name/{printf $2 " "}END{ print ""}' $file; done | sort -k 3 -n -r`
	cmd := exec.Command(
		"/bin/bash",
		"-c",
		memoryCmd,
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
			memory, err := strconv.Atoi(perProcess[2])
			if err != nil {
				log.Fatal(err)
			}
			if memory > 0 {
				p := ProcessMemory{Name: perProcess[0], PID: PID, MemoryKB: memory}
				p.MemoryPercent = float32(p.MemoryKB) / float32(m.MemoryTotalKB) * 100.0
				m.MemoryByProcess = append(m.MemoryByProcess, p)
			}
		}
	}
	ser, err := json.Marshal(m)
	fmt.Println(string(ser))
}

func (m *Memory) GetMemoryUsageTotal() {
	swapCmd := `cat /proc/meminfo | awk '/MemTotal|MemFree|Buffers|^Cached/{printf $2 " "}END{ print ""}'`
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
	memory := strings.Fields(out.String())
	m.MemoryTotalKB, err = strconv.Atoi(memory[0])
	if err != nil {
		log.Fatal(err)
	}
	free, err := strconv.Atoi(memory[1])
	if err != nil {
		log.Fatal(err)
	}
	buffers, err := strconv.Atoi(memory[2])
	if err != nil {
		log.Fatal(err)
	}
	cached, err := strconv.Atoi(memory[3])
	if err != nil {
		log.Fatal(err)
	}
	m.MemoryUsedKB = (m.MemoryTotalKB - free) - (buffers + cached)
	m.MemoryUsedPercent = float32(m.MemoryUsedKB) / float32(m.MemoryTotalKB) * 100.0
}
