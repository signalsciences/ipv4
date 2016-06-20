package ipv4

import (
	"net"
)

// IsIPv4 return true if input is a valid IPv4 address
func IsIPv4(dots string) bool {
	ip := net.ParseIP(dots)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// IsInternal returns true if the dotted IP address is in
//
// 0
// 10.0.0.0 - 10.255.255.255
// 127.0.0.0 - 127.255.255.255
// 172.16.0.0 - 172.31.255.255
// 192.168.0.0 - 192.168.255.255
//
// This might be better as "IsNotExternal" but that is gross.
//
func IsInternal(dots string) bool {
	ip := net.ParseIP(dots)
	if ip == nil {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 127 {
			return true
		}
		if ip4[0] == 10 {
			return true
		}
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] < 32 {
			return true
		}
		if ip4[0] == 0 && ip4[1] == 0 && ip4[2] == 0 && ip4[3] == 0 {
			return true
		}
	}

	// IpV6
	return false
}
