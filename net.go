package main

import (
	"net"
	"strings"
)

func IPv4(keys ...string) string {
	return findIP(keys, true)
}

func IPv6(keys ...string) string {
	return findIP(keys, false)
}

func findIP(keys []string, prefer4 bool) string {
	addrs, err := net.InterfaceAddrs()
	ips := filterByVersion(addrs, err, prefer4)

	for _, key := range keys {
		if nic, err := net.InterfaceByName(key); err == nil {
			addrs, err := nic.Addrs()
			if ips := filterByVersion(addrs, err, prefer4); len(ips) > 0 {
				return ips[0]
			}
		} else {
			for _, ip := range ips {
				if strings.HasPrefix(ip, key) {
					return ip
				}
			}
		}
	}
	return ""
}

func filterByVersion(addrs []net.Addr, err error, prefer4 bool) []string {
	var ips []string
	if err != nil {
		return ips
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err == nil && (ip.To4() != nil) == prefer4 {
			ips = append(ips, ip.String())
		}
	}

	return ips
}
