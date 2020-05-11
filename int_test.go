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

func TestToNetIP(t *testing.T) {
	n := ToNetIP(3232235777)
	ip := net.ParseIP("192.168.1.1")
	if !n.Equal(ip) {
		t.Errorf("ToNetIP(3232235777) = %v, want %v", n, ip)
	}
}

func TestFromDotsErrors(t *testing.T) {
	cases := []string{
		"....",
		"0",
		"0.",
		"0.0",
		"0.0.",
		"0.0.0",
		"0.0.0.",
		"0.0.0.0.",
		"0..",
		"0..0",
		"0.0.0..",
		"999.1.1.1",
		"1.1.1.11111",
	}
	for _, c := range cases {
		if _, err := FromDots(c); err == nil {
			t.Errorf("Case %q did not error", c)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	cases := []string{
		"0.0.0.0",
		"127.0.0.0",
		"128.0.0.0",
		"0.127.0.0",
		"0.128.0.0",
		"127.127.127.127",
		"128.128.128.128",
		"255.255.255.255",
	}

	for _, c := range cases {
		raw, err := FromDots(c)
		if err != nil {
			t.Errorf("Case %q didn't parse: %s", c, err)
			continue
		}
		if ToDots(raw) != c {
			t.Errorf("Case %q roundtrip produced %q", c, ToDots(raw))
		}
	}
}

var tempOut uint32

func BenchmarkFromDots(b *testing.B) {
	var val uint32
	for i := 0; i < b.N; i++ {
		val, _ = FromDots("255.129.128.127")
	}
	tempOut = val
}
