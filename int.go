package ipv4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// FromNetIP converts a IPv4 net.IP to uint32, error
func FromNetIP(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("not a IPv4 address")
	}
	return binary.BigEndian.Uint32(ip), nil
}

// FromDots converts a dotted IPv4 address to a uint32
// http://play.golang.org/p/T5B-6RExlj
// https://groups.google.com/forum/#!topic/golang-nuts/7sC28I57LRY
func FromDots(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("not a IPv4 address")
	}
	return binary.BigEndian.Uint32(ip), nil
}

// ToDots converts a uint32 to a IPv4 Dotted notation
func ToDots(val uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		val>>24,
		(val>>16)&0xFF,
		(val>>8)&0xFF,
		val&0xFF)
}
