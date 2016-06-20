package ipv4

import (
	"sort"
)

// Set defines a unique set of IPv4 addresses using a simple []uint32
type Set struct {
	Raw []uint32
}

// Len returns a length, part of the Sort.Interface
func (m Set) Len() int {
	return len(m.Raw)
}

// Swap is part of the Sort.Interface
func (m Set) Swap(i, j int) {
	m.Raw[i], m.Raw[j] = m.Raw[j], m.Raw[i]
}

// Less is part of the Sort.Interface
func (m Set) Less(i, j int) bool {
	return m.Raw[i] < m.Raw[j]
}

// Contains matches a dotted IPv4 address with an internal list
func (m Set) Contains(ipv4dots string) bool {
	if len(m.Raw) == 0 {
		return false
	}
	x, err := FromDots(ipv4dots)
	if err != nil {
		return false
	}
	i := sort.Search(len(m.Raw), func(i int) bool { return m.Raw[i] >= x })
	return (i < len(m.Raw) && m.Raw[i] == x)
}

// Add an element to the set
func (m *Set) Add(ipv4dots string) bool {
	x, err := FromDots(ipv4dots)
	if err != nil {
		return false
	}
	i := sort.Search(len(m.Raw), func(i int) bool { return m.Raw[i] >= x })
	if i == len(m.Raw) {
		// goes at end
		m.Raw = append(m.Raw, x)
		return true
	}
	if m.Raw[i] == x {
		// already exists
		return false
	}
	// splice
	// add one extra elemnt
	m.Raw = append(m.Raw, 0)
	// shift right
	copy(m.Raw[i+1:], m.Raw[i:])
	// insert new element in hole
	m.Raw[i] = x
	return true
}

// Valid return true if the internal storage is in sorted form and unique
func (m Set) Valid() bool {
	if len(m.Raw) == 0 {
		return true
	}
	last := m.Raw[0]
	for _, val := range m.Raw[1:] {
		if val <= last {
			return false
		}
		last = val
	}
	return true
}
