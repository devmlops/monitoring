package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"log"
)

func connections() string {
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
	return out.String()
}

func main() {
	conn := connections()
	fmt.Println(conn)
}

