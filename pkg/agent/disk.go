package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"log"
	"time"
)

type Disk struct {
	Time            time.Time `json:"time"`
	DiskTotalKB     uint64    `json:"disk_total_kb"`
	DiskUsedKB      uint64    `json:"disk_used_kb"`
	DiskUsedPercent float64   `json:"disk_used_percent"`
}

func (d *Disk) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	d.GetDiskUsage()
}

func (d *Disk) GetDiskUsage() {
	d.Time = time.Now().UTC()
	stat, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	d.DiskTotalKB = stat.Total / 1024
	d.DiskUsedKB = stat.Used / 1024
	d.DiskUsedPercent = stat.UsedPercent

	if Debug == true {
		ser, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(ser))
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(d)
	res, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s", server.IP, server.port, "disk"),
		"application/json; charset=utf-8",
		b,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
