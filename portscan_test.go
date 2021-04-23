package greg_test

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/cauefcr/greg"
)

func TestScanner(T *testing.T) {
	ports := []string{}
	i := 0
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:8666")
		if err != nil {
			T.Error(err)
		}
		conn, err := listener.Accept()
		if err != nil {
			T.Error(err)
		}
		conn.Write([]byte("the game"))
		conn.Close()
	}()
	ports = append(ports, fmt.Sprintf("800%v", i))
	for addr := range greg.PortScan([]string{"127.0.0.1"}, ports) {
		portEnd, err := strconv.ParseInt(strings.Split(addr, "8666")[1], 10, 32)
		if err != nil {
			T.Error(err)
		}
		if portEnd > 0 && portEnd < 100 {
			T.Errorf("listener closed")
		}
	}
}
