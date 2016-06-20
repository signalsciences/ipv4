package ipv4

import (
	"net"
	"testing"
)

func TestFromNetIPError(t *testing.T) {
	ip6 := net.ParseIP("2001:4860:0:2001::68")
	val, err := FromNetIP(ip6)
	if err == nil || val != 0 {
		t.Errorf("Expected error parsing IPv6 string")
	}
}

func TestFromDotsError(t *testing.T) {
	val, err := FromDots("2001:4860:0:2001::68")
	if err == nil || val != 0 {
		t.Errorf("Expected error parsing IPv6 string")
	}
}
