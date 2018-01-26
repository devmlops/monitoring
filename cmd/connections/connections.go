// WTF?????

package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i <= 100; i++ {
		go func() {
			defer wg.Done()
			_, err := net.Dial("tcp", "127.0.0.1:8090")
			time.Sleep(10 * time.Minute)
			if err != nil {
				fmt.Println("foo")
			}
		}()
	}
	wg.Wait()
}
