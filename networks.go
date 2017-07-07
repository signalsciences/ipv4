package ipv4

import (
	"net"
	"strings"
)

// IsIPv4 returns true if the input is either a dotted IPv4 address or if
// it's IPv4 dotted/cidr notation
func IsIPv4(s string) bool {
	var ipany net.IP
	if strings.IndexByte(s, '/') != -1 {
		ip, _ /*mask*/, err := net.ParseCIDR(s)
		if err != nil {
			return false
		}
		ipany = ip
	} else {
		ipany = net.ParseIP(s)
	}
	return ipany != nil && ipany.To4() != nil
}

// IsPrivate determines if a ip address is Not Public.
//
// Note in this case Private means "localhost, loopback, link local and private
// subnets".
//
func IsPrivate(ipdots string) bool {
	ip := net.ParseIP(ipdots)
	if ip != nil {
		ip4 := ip.To4()
		if ip4 != nil {
			switch {
			case ip4[0] == 127:
				return true
			// 10-net Class A 10.0.0.0/8
			case ip4[0] == 10:
				return true
			// 192.168.0.0/16
			case ip4[0] == 192 && ip4[1] == 168:
				return true
			// 172.16.0.0/12
			case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
				return true
			// link local  169.254.0.0/16
			case ip4[0] == 169 && ip4[1] == 254:
				return true
			}
		}
	}

	// sometimes we get stuff like localhost:2131231 (some port number)
	if strings.HasPrefix(ipdots, "localhost") || strings.HasPrefix(ipdots, "127.0.0.1:") {
		return true
	}

	return false
}
