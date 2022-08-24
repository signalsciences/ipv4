package ipv4

import (
	"errors"
	"net"
)

// ErrBadIP is a generic error that an IP address could not be parsed
var ErrBadIP = errors.New("Bad IP address")

// FromNetIP converts a IPv4 net.IP to uint32, error
func FromNetIP(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, ErrBadIP
	}
	return uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24, nil
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
// About 10x faster than doing something with fmt.Sprintf
// one allocation per call.
//
// Based on golang's net/IP.String()
// https://golang.org/src/net/ip.go?s=7645:7673#L281
func ToDots(p4 uint32) string {
	const maxIPv4StringLen = len("255.255.255.255")
	b := make([]byte, maxIPv4StringLen)

	n := ubtoa(b, 0, byte(p4>>24))
	b[n] = '.'
	n++

	n += ubtoa(b, n, byte((p4>>16)&0xFF))
	b[n] = '.'
	n++

	n += ubtoa(b, n, byte((p4>>8)&0xFF))
	b[n] = '.'
	n++

	n += ubtoa(b, n, byte(p4&0xFF))
	return string(b[:n])
}

// from
// https://golang.org/src/net/ip.go?s=7645:7673#L281
//
// ubtoa encodes the string form of the integer v to dst[start:] and
// returns the number of bytes written to dst. The caller must ensure
// that dst has sufficient length.
func ubtoa(dst []byte, start int, v byte) int {
	if v < 10 {
		dst[start] = v + '0'
		return 1
	} else if v < 100 {
		dst[start+1] = v%10 + '0'
		dst[start] = v/10 + '0'
		return 2
	}

	dst[start+2] = v%10 + '0'
	dst[start+1] = (v/10)%10 + '0'
	dst[start] = v/100 + '0'
	return 3
}

// SortUniqueUint32 sorts, and dedups a slice of uint32 (maybe representing
// binary representation of IPv4 address
//
// sorting and uniqueness is done in place
func SortUniqueUint32(in []uint32) {
	// reuse Set (which is a []unit32 anyways) implimentation
	set := Set(in)
	set.sort()
}
