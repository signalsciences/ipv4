package ipv4

import "testing"

func TestAdd(t *testing.T) {
	s := Set{}
	if s.Len() != 0 {
		t.Errorf("Size is not 0: %d", s.Len())
	}
	if s.Valid() == false {
		t.Errorf("invalid 0")
	}

	if s.Contains("12.12.12.12") == true {
		t.Errorf("Contains on empty set returned true")
	}

	if s.Add("12.12.12.12") == false {
		t.Errorf("Unable to insert")
	}
	if s.Len() != 1 {
		t.Errorf("Size is not 1: %d", s.Len())
	}
	if s.Valid() == false {
		t.Errorf("invalid 1")
	}
	if s.Add("12.12.12.12") == true {
		t.Errorf("Inserted duplicate")
	}

	if s.Contains("12.12.12.12") == false {
		t.Errorf("Unable to find 12.12.12.12")
	}

	if s.Add("1.1.1.1") == false {
		t.Errorf("Unable to insert 1")
	}
	if s.Len() != 2 {
		t.Errorf("Size is not 2: %d", s.Len())
	}
	if s.Valid() == false {
		t.Errorf("invalid 2")
	}
	if s.Contains("12.12.12.12") == false {
		t.Errorf("Unable to find 12.12.12.12")
	}
	if s.Contains("1.1.1.1") == false {
		t.Errorf("Unable to find 1.1.1.1")
	}
	if s.Add("6.6.6.6") == false {
		t.Errorf("Unable to insert 1")
	}
	if s.Len() != 3 {
		t.Errorf("Size is not 3: %d", s.Len())
	}
	if s.Valid() == false {
		t.Errorf("invalid 3")
	}
	if s.Contains("12.12.12.12") == false {
		t.Errorf("Unable to find 12.12.12.12")
	}
	if s.Contains("1.1.1.1") == false {
		t.Errorf("Unable to find 1.1.1.1")
	}
	if s.Contains("6.6.6.6") == false {
		t.Errorf("Unable to find 6.6.6.6")
	}

	if s.Add("junk") == true {
		t.Errorf("Accepted junk input")
	}
	if s.Contains("junk") == true {
		t.Errorf("Accepted junk input")
	}

	// intentionally have data be unsorted
	s = append(s, uint32(0))
	if s.Valid() {
		t.Errorf("have s.Valid() == true, want false")
	}

	s.sort()
	if !s.Valid() {
		t.Errorf("sorting did not work")
	}
}

func TestAddAll(t *testing.T) {
	s := Set{}
	s.AddAll([]string{"127.0.0.1", "10.0.0.1"})
	if !s.Valid() {
		t.Fatalf("AddAll is not valid")
	}
	out := s.ToDots()
	if out[0] != "10.0.0.1" || out[1] != "127.0.0.1" {
		t.Fatalf("Dump failed: %v", out)
	}
}
