package net

import (
	"fmt"
	"net"
	"strings"
)

func LocalIPAddr() (ipAddr string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addrs {
		addrStr := addr.String()
		if strings.Contains(addrStr, "127.0.0.1") ||
			strings.Contains(addrStr, "::1") ||
			strings.Contains(addrStr, "/64") {
			continue
		}

		addrStr = strings.TrimRight(addrStr, "1234567890")
		addrStr = strings.TrimRight(addrStr, "/")
		if len(addrStr) < 1 {
			continue
		}

		ipAddr = addrStr
		return
	}

	err = fmt.Errorf("no valid local ip address found")
	return
}
