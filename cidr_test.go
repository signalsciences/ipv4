package ipv4

import (
	"fmt"
	"log"
	"testing"
)

func TestCIDR2Range(t *testing.T) {
	tests := []struct {
		cidr  string
		left  string
		right string
		err   error
	}{
		{"10.0.1.32/27", "10.0.1.32", "10.0.1.63", nil},
		{"199.27.72.0/24", "199.27.72.0", "199.27.72.255", nil},
		{"192.0.2.100/24", "192.0.2.0", "192.0.2.255", nil},
		{"199.27.72.0/21", "199.27.72.0", "199.27.79.255", nil},
		{"192.168.100.0/22", "192.168.100.0", "192.168.103.255", nil},
		{"192.0.2.100/22", "192.0.0.0", "192.0.3.255", nil},
		{"192.0.0.100/16", "192.0.0.0", "192.0.255.255", nil},
		{"192.0.0.0/16", "192.0.0.0", "192.0.255.255", nil},
		{"169.254.0.0/16", "169.254.0.0", "169.254.255.255", nil},
		{"169.254.9.9/16", "169.254.0.0", "169.254.255.255", nil},
		{"172.16.0.0/12", "172.16.0.0", "172.31.255.255", nil},
		{"192.0.0.0/8", "192.0.0.0", "192.255.255.255", nil},
		{"192.0.2.0/8", "192.0.0.0", "192.255.255.255", nil},
		{"2001:db8:a0b:12f0::1/32", "", "", ErrBadIP},
	}
	for pos, tt := range tests {
		left, right, err := CIDR2Range(tt.cidr)
		if left != tt.left || right != tt.right || err != tt.err {
			t.Errorf("%d: %s got [%s, %s], want [%s, %s]", pos, tt.cidr, left, right, tt.left, tt.right)
		}
	}
}

func TestRange2CIDRs(t *testing.T) {
	tests := []struct {
		left   string
		right  string
		ranges []string
	}{
		{"127.0.0.1", "127.0.0.0", nil},
		{"junk", "127.0.0.0", nil},
		{"127.0.0.1", "junk", nil},
		{"0.0.0.0", "255.255.255.255", []string{"0.0.0.0/0"}},
		{"127.0.0.1", "127.0.0.18", []string{"127.0.0.1/32", "127.0.0.2/31", "127.0.0.4/30", "127.0.0.8/29", "127.0.0.16/31", "127.0.0.18/32"}},
		{"127.0.0.1", "127.0.0.1", []string{"127.0.0.1/32"}},
	}
	for _, tt := range tests {
		val := Range2CIDRs(tt.left, tt.right)
		if tt.ranges == nil && val != nil {
			t.Errorf("mismatch")
		} else {
			want := fmt.Sprintf("%v", tt.ranges)
			got := fmt.Sprintf("%v", val)
			if got != want {
				t.Errorf("Range2CIDRs: got %s, want %s", got, want)
			}
		}
	}
}

func ExampleCIDR2Range() {
	left, right, err := CIDR2Range("199.27.72.0/21")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(left, right)
	// Output: 199.27.72.0 199.27.79.255
}
func ExampleRange2CIDRs() {
	fmt.Println(Range2CIDRs("127.0.0.0", "127.0.0.255"))
	// Output: [127.0.0.0/24]
}
