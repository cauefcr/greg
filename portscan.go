package greg

import (
	"math"
	"net"
	"sync"
	"time"
)

// Nicelly parallel self contained function, returns a string chan which is fed with the open ip:port pairs as they are identified
func PortScan(ips, ports []string) chan string {
	c := make(chan string, 2048)
	done := make(chan string, 1024)
	// how to sync the workers
	wg := sync.WaitGroup{}
	// feeding the workers
	wg.Add(1)
	go func(c, done chan string, wg *sync.WaitGroup) {
		for _, ip := range ips {
			for _, p := range ports {
				c <- ip + ":" + p
			}
		}
		close(c)
		wg.Done()
	}(c, done, &wg)
	// start workers
	for i := 0; i < int(math.Sqrt(float64(len(ips)*len(ports)))); i++ {
		wg.Add(1)
		go func(c, done chan string, wg *sync.WaitGroup) {
			for addr := range c {
				// make the magic happen
				conn, err := net.DialTimeout("tcp", addr, time.Second)
				if err != nil {
					continue
				}
				done <- addr
				conn.Close()
			}
			wg.Done()
		}(c, done, &wg)
	}
	go func(wg *sync.WaitGroup, done chan string) {
		wg.Wait()
		close(done)
	}(&wg, done)
	return done
}
