package ipv4

import (
	"sort"
)

// Set defines a unique set of IPv4 addresses using a simple []uint32
//
// Since Set is a alias of []unit32, one can use
// `make(Set, length, capacity)` or use the NewSet constructor
type Set []uint32

// NewSet creates a Set with a given initial capacity.
func NewSet(capacity int) Set {
	return make(Set, 0, capacity)
}

// Len returns a length, part of the Sort.Interface
func (m Set) Len() int {
	return len(m)
}

// Swap is part of the Sort.Interface
func (m Set) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less is part of the Sort.Interface
func (m Set) Less(i, j int) bool {
	return m[i] < m[j]
}

// Contains matches a dotted IPv4 address with an internal list
func (m Set) Contains(ipv4dots string) bool {
	x, err := FromDots(ipv4dots)
	if err != nil {
		return false
	}
	i := sort.Search(len(m), func(i int) bool { return m[i] >= x })
	return i < len(m) && m[i] == x
}

// Add an element to the set.
func (m *Set) Add(ipv4dots string) bool {
	x, err := FromDots(ipv4dots)
	if err != nil {
		return false
	}
	orig := *m
	i := sort.Search(len(orig), func(i int) bool { return orig[i] >= x })
	if i == len(orig) {
		// goes at end
		*m = append(orig, x)
		return true
	}
	if orig[i] == x {
		// already exists
		return false
	}
	// splice
	// add one extra elemnt
	orig = append(orig, 0)
	// shift right
	copy(orig[i+1:], orig[i:])
	// insert new element in hole
	orig[i] = x
	*m = orig
	return true
}

// AddAll adds many IPv4 at once
func (m *Set) AddAll(ipv4dots []string) bool {
	in := *m
	for _, val := range ipv4dots {
		if bval, err := FromDots(val); err == nil {
			in = append(in, bval)
		}
	}
	*m = in
	m.sort()
	return true
}

// Valid return true if the internal storage is in sorted form and unique
func (m Set) Valid() bool {
	if len(m) == 0 {
		return true
	}
	last := m[0]
	for _, val := range m[1:] {
		if val <= last {
			return false
		}
		last = val
	}
	return true
}

// ToDots returns the IP set as a list of dotted-notation strings
func (m Set) ToDots() []string {
	out := make([]string, len(m))
	for i, val := range m {
		out[i] = ToDots(val)
	}
	return out
}

func (m *Set) sort() {
	in := *m
	sort.Sort(in)

	// inplace de-dup, uniqueness
	// https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	in = in[:j+1]
}
