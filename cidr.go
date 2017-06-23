package ipv4

import (
	"fmt"
	"net"
)

// CIDR2Range converts a CIDR to a dotted IP address pair, or empty strings and error
// Generic.. does not care if ipv4 or ipv6
func CIDR2Range(c string) (string, string, error) {
	left, ipnet, err := net.ParseCIDR(c)
	if err != nil {
		return "", "", err
	}
	left4 := left.To4()
	if left4 == nil {
		return "", "", nil
	}
	right := net.IPv4(0, 0, 0, 0).To4()
	right[0] = left4[0] | ^ipnet.Mask[0]
	right[1] = left4[1] | ^ipnet.Mask[1]
	right[2] = left4[2] | ^ipnet.Mask[2]
	right[3] = left4[3] | ^ipnet.Mask[3]

	return left4.String(), right.To4().String(), nil
}

// Range2CIDRs take a pair of IpV4 addresses in dotted form, and return a list of
// IPV4 CIDR ranges.  (or nil if invalid input)
func Range2CIDRs(dots1, dots2 string) (r []string) {
	a1, err := FromDots(dots1)
	if err != nil {
		return nil
	}
	a2, err := FromDots(dots2)
	if err != nil {
		return nil
	}
	if a1 > a2 {
		return nil
	}
	for a1 <= a2 {
		var l, first, last uint32
		for l = 32; l >= 0; l-- {
			m := (uint32(1) << l) - 1
			first = a1 & ^m
			last = a1 + m
			if a1 == first && last <= a2 {
				break
			}
		}
		r = append(r, fmt.Sprintf("%s/%d", ToDots(a1), 32-l))
		a1 = last
		if a1 == uint32(0xFFFFFFFF) {
			break
		}
		a1++
	}
	return r
}
