package agent

import (
	"net/http"
	"strconv"
	"encoding/json"
	"os/exec"
	"bytes"
	"time"
	"log"
	"fmt"
	"strings"
	"sync"
)

type Network struct {
	Time            time.Time    `json:"time"`
	Connections     int          `json:"connections"`
	ConnectionsByIP []Connection `json:"connections_by_ip"`
}

type Connection struct {
	IPAddress string    `json:"ip_address"`
	Number    int       `json:"number"`
}

//type Report interface {
//	Results()
//}

func (n *Network) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	n.GetActiveConnections()
}

func (n *Network) GetActiveConnections() {
	n.Time = time.Now().UTC()
	netstatCmd := "netstat -tn 2>/dev/null | tail -n +3 | awk '{print $5}' | cut -d: -f1 | sort | uniq -c | sort -nr"
	cmd := exec.Command(
		"/bin/bash",
		"-c",
		netstatCmd,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		perIP := strings.Fields(line)
		if len(perIP) > 0 {
			number, err := strconv.Atoi(perIP[0])
			if err != nil {
				log.Fatal(err)
			}
			c := Connection{IPAddress: perIP[1], Number: number}
			n.ConnectionsByIP = append(n.ConnectionsByIP, c)
			n.Connections += number
		}
	}
	ser, err := json.Marshal(n)
	fmt.Println(string(ser))
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(n)
	res, _ := http.Post("http://127.0.0.1:8080", "application/json; charset=utf-8", b)
	fmt.Println(res)
}

