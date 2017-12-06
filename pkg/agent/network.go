package agent

import (
	"sort"
	//	"log"
	"fmt"
	"sync"
	"time"
	//	"bytes"
	//	"net/http"
	"encoding/json"
	"github.com/shirou/gopsutil/net"
)

type Network struct {
	Time            time.Time    `json:"time"`
	Connections     int          `json:"connections"`
	ConnectionsByIP []Connection `json:"connections_by_ip"`
}

type Connection struct {
	IPAddress string `json:"ip_address"`
	Number    int    `json:"number"`
}

func (n *Network) RunJob(wg *sync.WaitGroup) {
	defer wg.Done()
	n.GetActiveConnections()
}

func (n *Network) GetActiveConnections() {
	n.Time = time.Now().UTC()

	cs, _ := net.Connections("tcp")

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
	nn := map[int][]string{}
	var a []int
	for k, v := range freq {
		nn[v] = append(nn[v], k)
	}
	for k := range nn {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range nn[k] {
			fmt.Printf("%s: %d\n", s, k)
			c := Connection{IPAddress: s, Number: k}
			n.ConnectionsByIP = append(n.ConnectionsByIP, c)
			n.Connections += k
		}
	}
	ser, _ := json.Marshal(n)
	fmt.Println(string(ser))
}
