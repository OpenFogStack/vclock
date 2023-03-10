// Package vclock implements vector clocks. Vector clocks are a way of capturing partial ordering of events in
// a distributed system. For more information see http://en.wikipedia.org/wiki/Vector_clock
//
// Create a new vector clock with the New function:
//
//	vc := vclock.New()
//
// Set the clock value for a given id:
//
//	vc.Set("A", 1)
//
// Tick the clock for a given id (increment the value by 1):
//
//	vc.Tick("A")
//
// Merge two vector clocks together:
//
//	vc.Merge(other)
//
// Compare two vector clocks:
//
//	vc.Compare(other, vclock.Equal)
//
// Get the relative ordering of two vector clocks:
//
//	vc.Order(other)
package vclock

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sort"
)

// Condition constants define how to compare a vector clock against another,
// and may be ORed together when being provided to the Compare method.
// Use the Equal, Ancestor, Descendant, and Concurrent constants to
// compare vector clocks.
type Condition int

const (
	Equal Condition = 1 << iota
	Ancestor
	Descendant
	Concurrent
)

// VClock is the type of a vector clock. A VClock is a map of process ids (string)
// to clock values (uint64). Note that this implies that the zero clock value is 0.
type VClock map[string]uint64

// FindTicks returns the clock value for a given id or false if the id is not found.
// This is a convenience function for accessing the clock value of a given id, you
// can also access the clock value directly using the map syntax, e.g.: vc["A"].
func (vc VClock) FindTicks(id string) (uint64, bool) {
	ticks, ok := vc[id]
	return ticks, ok
}

// New returns a new, empty vector clock.
func New() VClock {
	return VClock{}
}

// Copy returns a deep copy of a vector clock.
func (vc VClock) Copy() VClock {
	cp := make(map[string]uint64, len(vc))
	for key, value := range vc {
		cp[key] = value
	}
	return cp
}

// CopyFromMap creates a shallow copy of a map that is compatible with the VClock type.
func (vc VClock) CopyFromMap(otherMap map[string]uint64) VClock {
	return otherMap
}

// GetMap returns the underlying map of the vector clock.
func (vc VClock) GetMap() map[string]uint64 {
	return vc
}

// Set sets the clock value of the given process id to the given value.
func (vc VClock) Set(id string, ticks uint64) {
	vc[id] = ticks
}

// Tick increments the clock value of the given process id by 1.
// If the process id is not found in the clock, it is added with a value of 1.
func (vc VClock) Tick(id string) {
	vc[id] = vc[id] + 1
}

// LastUpdate returns the smallest clock value in the vector clock.
func (vc VClock) LastUpdate() (last uint64) {
	for key := range vc {
		if vc[key] > last {
			last = vc[key]
		}
	}
	return last
}

// Merge takes the maximum of all clock values in other and updates the
// values of the callee. If the callee does not contain a given id, it is
// added to the callee with the value from other.
// Merge updates the callee vector clock in place.
func (vc VClock) Merge(other VClock) {
	for id := range other {
		if vc[id] < other[id] {
			vc[id] = other[id]
		}
	}
}

// Bytes returns an encoded vector clock using the gob package.
func (vc VClock) Bytes() []byte {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(vc)
	if err != nil {
		log.Fatal("Vector Clock Encode:", err)
	}
	return b.Bytes()
}

// FromBytes decodes a vector clock from a byte slice using the gob package.
func FromBytes(data []byte) (vc VClock, err error) {
	b := new(bytes.Buffer)
	b.Write(data)
	clock := New()
	dec := gob.NewDecoder(b)
	err = dec.Decode(&clock)
	return clock, err
}

// PrintVC prints the callees vector clock to stdout.
func (vc VClock) PrintVC() {
	fmt.Println(vc.ReturnVCString())
}

// ReturnVCString returns a deterministic string encoding of a vector clock.
func (vc VClock) ReturnVCString() string {
	//sort
	ids := make([]string, len(vc))
	i := 0
	for id := range vc {
		ids[i] = id
		i++
	}

	sort.Strings(ids)

	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i := range ids {
		buffer.WriteString(fmt.Sprintf("\"%s\":%d", ids[i], vc[ids[i]]))
		if i+1 < len(ids) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

// Order determines the relationship between two clocks. It returns
// Ancestor if the given clock other is an ancestor of the callee vc,
// Descendant if other is a descendant of vc, Equal if other and vc
// are equal, and Concurrent if other and vc are concurrent.
//
// Two important notes about this implementation:
//  1. The return value is a constant, it is not ORed together. This means
//     that if you want to compare the output, you can use the == operator.
//     Note that we recommend using the Compare method instead.
//  2. If the clocks are equal, the return value is Equal. This is different
//     from the original vector clock implementation, which returned Concurrent
//     AND Equal.
func (vc VClock) Order(other VClock) Condition {
	// This code is adapted from the Voldemort implementation of vector clocks,
	// see https://github.com/voldemort/voldemort/blob/master/src/java/voldemort/versioning/VectorClockUtils.java
	//
	// The original code is licensed under the Apache License, Version 2.0
	// http://www.apache.org/licenses/LICENSE-2.0
	// Copyright 2008-2013 LinkedIn, Inc.

	vcBigger := false
	otherBigger := false

	// we're finding the entries that both clocks have in common
	commonEntries := make(map[string]struct{})

	for id := range vc {
		if _, ok := other[id]; ok {
			commonEntries[id] = struct{}{}
		}
	}

	if len(vc) > len(commonEntries) {
		vcBigger = true
	}
	if len(other) > len(commonEntries) {
		otherBigger = true
	}

	for id := range commonEntries {
		if vcBigger && otherBigger {
			break
		}
		vcVersion := vc[id]
		otherVersion := other[id]

		if vcVersion > otherVersion {
			vcBigger = true
		} else if vcVersion < otherVersion {
			otherBigger = true
		}
	}

	if !vcBigger && !otherBigger {
		return Equal
	}

	if vcBigger && !otherBigger {
		return Ancestor
	}

	if !vcBigger && otherBigger {
		return Descendant
	}

	return Concurrent
}

// Compare takes another clock ("other") and determines if it is Equal, an
// Ancestor, Descendant, or Concurrent with the callees ("vc") clock.
// The condition is specified by the cond parameter, which may be ORed.
// For example, to check if two clocks are concurrent or descendants, you would
// call Compare(other, Concurrent|Descendant). If the condition is met, true
// is returned, otherwise false is returned.
func (vc VClock) Compare(other VClock, cond Condition) bool {
	return vc.Order(other)&cond != 0
}

// CompareOld takes another clock and determines if it is Equal, an
// Ancestor, Descendant, or Concurrent with the callees clock.
// Deprecated: This is the original implementation of Compare, which is now
// deprecated.  It is left here for reference.
func (vc VClock) CompareOld(other VClock, cond Condition) bool {
	var otherIs Condition
	// Preliminary qualification based on length
	if len(vc) > len(other) {
		if cond&(Ancestor|Concurrent) == 0 {
			return false
		}
		otherIs = Ancestor
	} else if len(vc) < len(other) {
		if cond&(Descendant|Concurrent) == 0 {
			return false
		}
		otherIs = Descendant
	} else {
		otherIs = Equal
	}

	//Compare matching items
	for id := range other {
		if _, found := vc[id]; found {
			if other[id] > vc[id] {
				switch otherIs {
				case Equal:
					otherIs = Descendant
					break
				case Ancestor:
					return cond&Concurrent != 0
				}
			} else if other[id] < vc[id] {
				switch otherIs {
				case Equal:
					otherIs = Ancestor
					break
				case Descendant:
					return cond&Concurrent != 0
				}
			}
		} else {
			if otherIs == Equal {
				return cond&Concurrent != 0
			} else if (len(other) - len(vc) - 1) < 0 {
				return cond&Concurrent != 0
			}
		}
	}

	//Equal clocks are concurrent
	if otherIs == Equal && cond == Concurrent {
		cond = Equal
	}
	return cond&otherIs != 0
}
