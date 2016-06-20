package ipv4

import (
	"fmt"
	"testing"
)

func TestCidrToRange(t *testing.T) {
	tests := []struct {
		cidr  string
		left  string
		right string
		err   error
	}{
		{"199.27.72.0/21", "199.27.72.0", "199.27.79.255", nil},
	}
	for pos, tt := range tests {
		left, right, err := CIDR2Range(tt.cidr)
		if left != tt.left || right != tt.right || err != tt.err {
			t.Errorf("%d: %s expected [%s, %s] got [%s, %s]", pos, tt.cidr, tt.left, tt.right, left, right)
		}
	}
}

func TestCidr(t *testing.T) {
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
			expected := fmt.Sprintf("%v", tt.ranges)
			actual := fmt.Sprintf("%v", val)
			if expected != actual {
				t.Errorf("Expected %s, Got %s", expected, actual)
			}
		}
	}
}
