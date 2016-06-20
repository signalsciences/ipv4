package ipv4

import (
	"net"
	"testing"
)

func TestInternal(t *testing.T) {
	cases := []struct{
		ip string,
		internal bool,
	}{
		{ "junk", true },
		{ "", true },
		{ "0.0.0.0", true},
		{ "127.0.0.1", true},
		{ "9999.9999.9999.999", true},
		{ "10.0.0.0", true},
	}
	for  d, tt := range cases {

	}	

}

