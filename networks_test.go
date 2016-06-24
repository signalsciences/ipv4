package ipv4

import "testing"

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"10.0.0.0", true},
		{"10.0.0.0/8", true},
		{"2001:4860:0:2001::68", false},
		{"2001:DB8::/48", false},
		{"false/48", false},
		{"false", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsIPv4(tt.ip); got != tt.want {
			t.Errorf("IsIPv4(%q) = %t, want %t", tt.ip, got, tt.want)
		}
	}
}

func TestIsPrivate(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		// Private IPs
		{"10.0.0.0", true},
		{"10.255.255.255", true},
		{"192.168.0.0", true},
		{"192.168.255.255", true},
		{"172.16.0.0", true},
		{"172.31.255.255", true},
		{"169.254.0.0", true},
		{"169.254.255.255", true},
		{"127.0.0.1", true},

		// Public IPs
		{"9.255.255.255", false},
		{"11.0.0.0", false},
		{"192.167.255.255", false},
		{"192.169.0.0", false},
		{"172.15.255.255", false},
		{"172.32.0.0", false},
		{"169.253.255.255", false},
		{"169.255.0.0", false},

		// Oddballs
		{"localhost", true},
		{"localhost:12312", true},
		{"127.0.0.1:12321", true},
	}

	for _, tt := range tests {
		if got := IsPrivate(tt.ip); got != tt.want {
			t.Errorf("IsPrivate(%q) = %t, want %t", tt.ip, got, tt.want)
		}
	}

}
