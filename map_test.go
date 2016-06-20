package ipv4

import (
	"testing"
)

func TestSettingRange(t *testing.T) {
	set := NewIntervalMap(100)

	// Smoke test to sort out major pass-by-val vs. pass-by-ref problems
	err := set.AddRange("20.0.0.0", "20.0.0.255", true)
	if err != nil || set.Len() != 1 {
		t.Fatalf("Adding failed: err=%v, len=%d", err, set.Len())
	}

	table := []struct {
		left  string
		right string
		ok    bool
	}{
		{"10.0.0.0", "11.0.0.0", false}, // class A
		{"1.0.0.0", "1.255.255.255", true},
		{"12.1.0.0", "12.0.0.0", false},
		{"2.0.0.0", "2.0.0.0", true},
		{"Busted", "2.0.0.0", false},
		{"2.0.0.0", "Busted", false},
	}

	for pos, test := range table {
		err := set.AddRange(test.left, test.right, true)
		if err == nil && test.ok == false {
			t.Errorf("test %d: Expected an error with input [%s, %s]", pos, test.left, test.right)
		} else if err != nil && test.ok == true {
			t.Errorf("test %d: Got an error with input [%s, %s]: %s", pos, test.left, test.right, err)
		}
		err = set.Valid()
		if err != nil {
			t.Errorf("Set state is invalid: %s", err)
		}
	}

	src := set.Go()
	if len(src) == 0 {
		t.Errorf("Output was empty")
	}
}

func TestSettingCIDR(t *testing.T) {
	set := NewIntervalMap(100)
	addtable := []struct {
		cidr string
		ok   bool
	}{
		{"Busted", false},
		{"2.0.0.0/", false},
		{"2.0.0.0/Busted", false},
		{"192.128.1.0/7", false},
		{"1.0.0.0/8", true},
		{"10.0.0.0/32", true},
		{"10.0.0.1/32", true},
		{"10.0.0.3", true},
	}
	for pos, test := range addtable {
		err := set.Add(test.cidr, pos)
		if err == nil && test.ok == false {
			t.Errorf("test %d: Expected an error with input %q", pos, test.cidr)
		} else if err != nil && test.ok == true {
			t.Errorf("test %d: Got an error with input %q: %s", pos, test.cidr, err)
		}
		err = set.Valid()
		if err != nil {
			t.Errorf("Set state is invalid: %s", err)
		}
	}
	if set.Len() != 4 {
		t.Errorf("expected len of 4, got %d", set.Len())
	}
	set.Add("10.0.0.4", 7)
	if set.Len() != 4 {
		t.Errorf("expected len of 4, got %d", set.Len())
	}
	//fmt.Printf("--->  INTERNAL TREE:\n %s\n", set)

	table := []struct {
		dots   string
		result interface{}
	}{
		{"0.0.0.0", nil},
		{"0.0.0.255", nil},
		{"1.0.0.0", 4},
		{"1.255.255.255", 4},
		{"2.0.0.0", nil},
		{"9.255.255.255", nil},
		{"10.0.0.0", 5},
		{"10.0.0.1", 6},
		{"10.0.0.2", nil},
		{"10.0.0.3", 7},
		{"10.0.0.4", 7},
		{"255.255.255.255", nil},
	}
	for pos, test := range table {
		val := set.Contains(test.dots)
		if val != test.result {
			t.Errorf("test %d: Contains(%q) is %v, expected %v", pos, test.dots, val, test.result)
		}
	}
}

/*
func BenchmarkLookup(b *testing.B) {

	for n := 0; n < b.N; n++ {
		Find("0.255.255.255")
	}
}
*/

// This test various way intervals can overlap with each other
// the code should consolidate them appropriately
func TestOverlaps(t *testing.T) {
	set := NewIntervalMap(100)

	table := []struct {
		left  string
		right string
		count int
	}{
		{"10.0.0.0", "10.0.1.127", 1},   // class A
		{"10.0.0.0", "10.0.1.127", 1},   // dup
		{"12.0.0.0", "12.1.0.0", 2},     // disjoint
		{"10.0.0.0", "10.0.1.132", 2},   // extension
		{"10.0.0.10", "10.0.1.131", 2},  // subset
		{"10.0.0.10", "10.0.1.132", 2},  // subset
		{"10.0.0.10", "10.0.1.135", 2},  // overlap
		{"10.0.0.135", "10.0.1.140", 2}, // continuation
	}
	for pos, test := range table {
		err := set.AddRange(test.left, test.right, pos)
		if err != nil {
			t.Fatalf("test %d: Error with input [%s, %s]: %s", pos, test.left, test.right, err)
		}
		if set.Len() != test.count {
			t.Errorf("Got bad count: expected %d got %d", test.count, set.Len())
		}
	}
}

// This test various way intervals can overlap with each other
// the code should consolidate them appropriately
func TestOverlaps2(t *testing.T) {
	set := NewIntervalMap(100)

	table := []string{
		"64.39.96.0/20",
		"64.39.106.0/20",
		"64.39.106.208",
	}
	for pos, test := range table {
		err := set.Add(test, pos)
		if err != nil {
			t.Fatalf("test %d: Error with input %s: %s", pos, test, err)
		}
		err = set.Valid()
		if err != nil {
			t.Errorf("Set state is invalid: %s", err)
		}
	}
}

func TestEmpty(t *testing.T) {
	set := NewIntervalMap(100)
	err := set.Valid()
	if err != nil {
		t.Errorf("Set state is invalid: %s", err)
	}
	result := set.Contains("127.0.0.1")

	if result != nil {
		t.Errorf("Empty set contained something!!")
	}

	result = set.Contains("junk")
	if result != nil {
		t.Errorf("Empty set contained something and parsed invalid input!!")
	}
}
