package main

import (
	"fmt"
	"greg"
	"math"
	"strings"
	"time"

	ranges "github.com/activeshadow/libminimega/ranges"
	"github.com/cauefcr/greg"
	"github.com/jessevdk/go-flags"
	cidr "github.com/nytr0gen/go-cidr"
)

// flag
type opt struct {
	Ports []string `short:"p" long:"ports" description:"Ports to scan" default:"[22-80]"`
}

func main() {
	start := time.Now()
	// parse flags
	opts := opt{}
	args, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	if len(args) == 0 {
		return
	}
	// parse rest of arguments as cidr if they have /
	realIPs := []string{}
	for _, ip := range args {
		if strings.Contains(ip, "/") {
			ips, err := cidr.List(ip)
			if err != nil {
				panic(err)
			}
			realIPs = append(realIPs, ips...)
		} else {
			realIPs = append(realIPs, ip)
		}
	}
	// expand port if they're ranges
	realPorts := []string{}
	for _, port := range opts.Ports {
		pts, err := ranges.SplitList(port)
		if err != nil {
			panic(err)
		}
		realPorts = append(realPorts, pts...)
	}
	// do the work
	for addr := range greg.PortScan(realIPs, realPorts) {
		fmt.Printf("%v\topen\n", addr)
	}
	fmt.Println("Took us", time.Since(start), "to get here, with", int(math.Sqrt(float64(len(realIPs)*len(realPorts)))), "coroutines")
}

// // Nicelly parallel self contained function, returns a string chan which is fed with the open ip:port pairs as they are found
// func PortScan(ips, ports []string) chan string {
// 	c := make(chan string, 2048)
// 	done := make(chan string, 1024)
// 	// how to sync the workers
// 	wg := sync.WaitGroup{}
// 	// feeding the workers
// 	wg.Add(1)
// 	go func(c, done chan string, wg *sync.WaitGroup) {
// 		for _, ip := range ips {
// 			for _, p := range ports {
// 				c <- ip + ":" + p
// 			}
// 		}
// 		close(c)
// 		wg.Done()
// 	}(c, done, &wg)
// 	// start workers
// 	for i := 0; i < int(math.Sqrt(float64(len(ips)*len(ports)))); i++ {
// 		wg.Add(1)
// 		go func(c, done chan string, wg *sync.WaitGroup) {
// 			for addr := range c {
// 				// make the magic happen
// 				conn, err := net.DialTimeout("tcp", addr, time.Second)
// 				if err != nil {
// 					continue
// 				}
// 				done <- addr
// 				conn.Close()
// 			}
// 			wg.Done()
// 		}(c, done, &wg)
// 	}
// 	go func(wg *sync.WaitGroup, done chan string) {
// 		wg.Wait()
// 		close(done)
// 	}(&wg, done)
// 	return done
// }
