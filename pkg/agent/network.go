package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"log"
	"sort"
	"time"
)

type Network struct {
	Time            time.Time    `json:"time"`
	Connections     uint64       `json:"connections"`
	ConnectionsByIP []Connection `json:"connections_by_ip"`
	Server          Server        `json:"-"`
	Debug           bool          `json:"-"`
	Hostname        string        `json:"hostname"`
}

type Connection struct {
	IPAddress string `json:"ip_address"`
	Number    uint64 `json:"number"`
}

func (n *Network) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	n.GetActiveConnections()
}

func (n *Network) GetActiveConnections() {
	n.ConnectionsByIP = nil
	n.Time = time.Now().UTC()

	cs, err := net.Connections("tcp")
	if err != nil {
		log.Println(err)
	}

	freq := make(map[string]int)
	for _, conn := range cs {
		if (conn.Status == "ESTABLISHED") && (conn.Raddr.IP != "127.0.0.1") {
			_, ok := freq[conn.Raddr.IP]
			if ok == true {
				freq[conn.Raddr.IP] += 1
			} else {
				freq[conn.Raddr.IP] = 1
			}
		}

	}
	reversed_freq := map[int][]string{}
	var numbers []int
	for key, val := range freq {
		reversed_freq[val] = append(reversed_freq[val], key)
	}
	for val := range reversed_freq {
		numbers = append(numbers, val)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	if len(numbers) > 5 {
		numbers = numbers[:5]
	}
	for _, number := range numbers {
		for _, s := range reversed_freq[number] {
			c := Connection{IPAddress: s, Number: uint64(number)}
			n.ConnectionsByIP = append(n.ConnectionsByIP, c)
			if len(n.ConnectionsByIP) > 4 {
				break
			}
			n.Connections += uint64(number)
		}
	}

	if n.Debug == true {
		ser, err := json.Marshal(n)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(ser))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(n)
	res, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s", n.Server.IP, n.Server.Port, "network"),
		"application/json; charset=utf-8",
		b,
	)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
}
