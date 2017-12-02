package agent

import (
	//"encoding/json"
	"os/exec"
	//"strconv"
	"bytes"
	"time"
	"log"
	//"fmt"
)

type Network struct {
	Time time.Time             `json:"time"`
	Connections []Connection   `json:"connections"`
}

type Connection struct {
	IPAddress string    `json:"ip_address"`
	Number    int       `json:"number"`
}

//type Report interface {
//	Results()
//}

func (n *Network) GetActiveConnections() {
	n.Time = time.Now().UTC()
	netstatCmd := "netstat -tn 2>/dev/null | tail -n +3 | awk '{print $5}' | cut -d: -f1 | sort | uniq -c | sort -nr | head"
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
	c := Connection{IPAddress: "8.8.8.8", Number: 24}
	n.Connections = append(n.Connections, c)
	c = Connection{IPAddress: "192.168.0.15", Number: 3}
	n.Connections = append(n.Connections, c)
	//serialized, err := json.Marshal(n)
}

