package ipv4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var ErrBadIP = errors.New("Bad IP address")

// FromNetIP converts a IPv4 net.IP to uint32, error
func FromNetIP(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("not a IPv4 address")
	}
	return binary.BigEndian.Uint32(ip), nil
}

// ToNetIP converts a uint32 to a net.IP (net.IPv4 actually)
func ToNetIP(val uint32) net.IP {
	return net.IPv4(byte(val>>24), byte(val>>16&0xFF),
		byte(val>>8)&0xFF, byte(val&0xFF))
}

// FromDots converts an IPv4 dotted string into a
// uint32 (big endian).
func FromDots(ipstr string) (uint32, error) {
	var out uint32
	var oct uint32
	var num int
	var hasOct bool

	for i := 0; i < len(ipstr); i++ {
		b := ipstr[i]
		switch {
		case b >= '0' && b <= '9':
			oct = oct*10 + uint32(b-'0')
			if oct > 255 {
				return 0, ErrBadIP
			}
			hasOct = true
		case b == '.':
			// test for ".." cases
			if !hasOct {
				return 0, ErrBadIP
			}
			out = (out << 8) | oct
			oct = 0
			hasOct = false
			num++

			// more than 3 dots
			if num > 3 {
				return 0, ErrBadIP
			}
		default:
			return 0, ErrBadIP
		}
	}

	// num != 3 --> must have 3 dots
	// !hasOct --> input ended in a dot or other char
	if num != 3 || !hasOct {
		return 0, ErrBadIP
	}
	out = (out << 8) | oct
	return out, nil
}

// ToDots converts a uint32 to a IPv4 Dotted notation
//
// Impl is mostly for debugging and is not high performance
func ToDots(val uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		val>>24,
		(val>>16)&0xFF,
		(val>>8)&0xFF,
		val&0xFF)
}
