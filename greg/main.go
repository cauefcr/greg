package main

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	ranges "github.com/activeshadow/libminimega/ranges"
	greg "github.com/cauefcr/greg"
	"github.com/jessevdk/go-flags"
	// cidr "github.com/knownsec/Minitools-cidrgen"
)

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func cidrToHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

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
			ips, err := cidrToHosts(ip)
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
