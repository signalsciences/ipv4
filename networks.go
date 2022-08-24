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

// IsPrivate determines if an IP address is Not Public.
//
// Note in this case Private means "localhost, loopback, link local and private
// subnets".
func IsPrivate(ipdots string) bool {
	ip := net.ParseIP(ipdots)
	if ip == nil {
		// sometimes we get stuff like localhost:2131231 (some port number)
		if strings.HasPrefix(ipdots, "localhost") || strings.HasPrefix(ipdots, "127.0.0.1:") {
			return true
		}
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() || ip.IsPrivate() {
		return true
	}
	return false
}
