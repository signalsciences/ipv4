package ipv4

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strings"
)

// Interval is a closed interval [a,b] (inclusive), with a value
type Interval struct {
	Left  uint32
	Right uint32
	Value interface{}
}

func (i Interval) String() string {
	return fmt.Sprintf("[%s, %s]=%v", ToDots(i.Left), ToDots(i.Right), i.Value)
}

// IntervalList is listing of Interval types, most for use in sorting
type IntervalList []Interval

func (ipset IntervalList) Len() int {
	return len(ipset)
}
func (ipset IntervalList) Less(i, j int) bool {
	return ipset[i].Left < ipset[j].Left
}
func (ipset IntervalList) Swap(i, j int) {
	ipset[i], ipset[j] = ipset[j], ipset[i]
}

func (ipset IntervalList) String() string {
	buf := bytes.Buffer{}
	for pos, val := range ipset {
		buf.WriteString(fmt.Sprintf("%d: %s\n", pos, val))
	}
	return buf.String()
}

// IntervalMap is set of disjoint intervals
type IntervalMap struct {
	Intervals IntervalList
}

// NewIntervalMap creates a new set
func NewIntervalMap(capacity int) *IntervalMap {
	return &IntervalMap{
		Intervals: make([]Interval, 0, capacity),
	}
}

// Go emits a source code representation of the data
func (ipset IntervalMap) Go() string {
	buf := bytes.Buffer{}
	buf.WriteString("ipv4.IntervalMap{\n")
	buf.WriteString("Intervals: []ipv4.Interval{\n")
	for _, val := range ipset.Intervals {
		buf.WriteString(fmt.Sprintf("%#v, // [%s, %s]\n", val, ToDots(val.Left), ToDots(val.Right)))
	}
	buf.WriteString("},\n}")
	return buf.String()
}

func (ipset IntervalMap) String() string {
	return ipset.Intervals.String()
}

func (ipset *IntervalMap) add(left, right uint32, value interface{}) error {
	if left > right {
		return fmt.Errorf("left %s > right %s",
			ToDots(left), ToDots(right))
	}
	if right-left >= (uint32(1) << 24) {
		return fmt.Errorf("Interval too large: [%s,%s]",
			ToDots(left), ToDots(right))
	}

	ipset.Intervals = append(ipset.Intervals, Interval{left, right, value})
	if len(ipset.Intervals) == 1 {
		return nil
	}

	sort.Sort(ipset.Intervals)

	// [-----]
	// [-----]     (duplicate)
	// [--------]  (extension)
	//   [--------] (extension)
	//        [------] (continuation)
	//   [--]         (subset)
	//          [ disjoint ]
	newset := make([]Interval, 1, len(ipset.Intervals)+1)

	last := ipset.Intervals[0]
	newset[0] = last
	for _, val := range ipset.Intervals[1:] {
		// equals or subset.. skip
		if val.Left >= last.Left && val.Right <= last.Right {
			continue
		}

		if val.Left >= last.Left && val.Left <= last.Right {
			last.Right = right
			newset[len(newset)-1] = last
			continue
		}
		// extend previous interval if the end of one is the start of another
		// AND the values are the same
		if val.Left == last.Right+1 && val.Value == last.Value {
			last.Right = val.Right
			newset[len(newset)-1] = last
			continue
		}
		last = val
		newset = append(newset, last)
	}

	ipset.Intervals = newset
	return ipset.Valid()
}

// Valid return error if internally invalid or nil if correct
func (ipset IntervalMap) Valid() error {
	last := Interval{}
	for pos, val := range ipset.Intervals {
		if val.Left > val.Right {
			return fmt.Errorf("left %s > right %s at pos %d",
				ToDots(val.Left), ToDots(val.Right), pos)
		}
		if val.Right-val.Left > (uint32(255) << 24) {
			return fmt.Errorf("Interval too large: [%s,%s]",
				ToDots(val.Left), ToDots(val.Right))
		}
		if pos > 0 {
			if val.Left <= last.Right || val.Right <= last.Right {
				return fmt.Errorf("Overlapping regions [%s,%s] vs. [%s,%s]",
					ToDots(last.Left), ToDots(last.Right),
					ToDots(val.Left), ToDots(val.Right))
			}
		}
		last = val
	}
	return nil
}

// Add inserts a single IP or a range based with CIDR notation
func (ipset *IntervalMap) Add(dots string, value interface{}) error {
	var left, right uint32
	var err error
	// Pure IP4 address
	if strings.IndexByte(dots, '/') == -1 {
		left, err = FromDots(dots)
		if err != nil {
			return fmt.Errorf("Unable to parse %q", dots)
		}
		right = left
	} else {
		// It's a CIDR
		leftip, cidrnet, err := net.ParseCIDR(dots)
		if err != nil {
			return err
		}
		ones, _ := cidrnet.Mask.Size()
		left, _ = FromNetIP(leftip)
		right = left + 1<<uint(32-ones) - 1
	}
	//fmt.Printf("ADDING [%s,%s]\n", ToDots(left), ToDots(right))
	return ipset.add(left, right, value)
}

// AddRange adds a range of IP addresses
func (ipset *IntervalMap) AddRange(dotsleft, dotsright string, value interface{}) error {
	left, err := FromDots(dotsleft)
	if err != nil {
		return fmt.Errorf("Unable to parse %q", dotsleft)
	}
	right, err := FromDots(dotsright)
	if err != nil {
		return fmt.Errorf("Unable to parse %q", dotsright)
	}
	return ipset.add(left, right, value)
}

// Len returns the number of intervals in the set
func (ipset IntervalMap) Len() int {
	return ipset.Intervals.Len()
}

// Contains returns true if the ip is in the set
//   error if set or input isn't valid
func (ipset IntervalMap) Contains(dots string) interface{} {
	val, err := FromDots(dots)
	if err != nil {
		return nil
	}
	ilen := ipset.Intervals.Len()
	if ilen == 0 {
		return nil
	}

	i := sort.Search(ilen, func(i int) bool {
		return ipset.Intervals[i].Left >= val
	})

	// if we overflowed, then check the last value
	//  we know i-1 > 0 (safe), since we checked for this above
	if i == ilen {
		i--
		if ipset.Intervals[i].Left <= val && val <= ipset.Intervals[i].Right {
			return ipset.Intervals[i].Value
		}
		return nil
	}

	// Did the IP match the exact start of an interval?
	if ipset.Intervals[i].Left == val {
		return ipset.Intervals[i].Value
	}

	// if we are at the start, then no match
	if i == 0 {
		return nil
	}

	// safe
	i--
	if ipset.Intervals[i].Left <= val && val <= ipset.Intervals[i].Right {
		return ipset.Intervals[i].Value
	}
	return nil
}
