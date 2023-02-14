package vclock

import (
	"fmt"
	"testing"
)

func TestBasicInit(t *testing.T) {
	n := New()
	n.Set("a", 2)
	n.Set("b", 1)
	na, isFounda := n.FindTicks("a")

	if !isFounda {
		t.Fatalf("Failed on finding ticks: %s", n.ReturnVCString())
	}

	if na != 2 {
		t.Fatalf("Tick value did not increment: %s", n.ReturnVCString())
	}

	n.Tick("b")

	na, isFounda = n.FindTicks("a")
	nb, isFoundb := n.FindTicks("b")

	if !isFounda || !isFoundb {
		t.Fatalf("Failed on finding ticks: %s", n.ReturnVCString())
	}

	if na != 2 || nb != 2 {
		t.Fatalf("Tick value did not increment: %s", n.ReturnVCString())
	}

}

func TestCopy(t *testing.T) {
	n := New()
	n.Set("a", 4)
	n.Set("b", 1)
	n.Set("c", 3)
	n.Set("d", 2)
	nc := n.Copy()

	an, _ := nc.FindTicks("a")
	bn, _ := nc.FindTicks("b")
	cn, _ := nc.FindTicks("c")
	dn, _ := nc.FindTicks("d")

	ao, _ := n.FindTicks("a")
	bo, _ := n.FindTicks("b")
	co, _ := n.FindTicks("c")
	do, _ := n.FindTicks("d")

	if an != ao || bn != bo || cn != co || dn != do {
		failComparison(t, "Copy not the same as the original new = %s , old = %s ", nc, n)
	} else if !n.Compare(nc, Equal) {
		failComparison(t, "Copy not the same as the original new = %s , old = %s ", n, nc)
	}
}

func TestCompareAndMerge(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 2)
	n1.Set("b", 1)
	n1.Set("c", 1)

	n2.Set("a", 1)
	n2.Set("b", 3)
	n2.Set("c", 1)

	n3 := n1.Copy()
	n3.Merge(n2)

	an, _ := n3.FindTicks("a")
	bn, _ := n3.FindTicks("b")
	cn, _ := n3.FindTicks("c")

	nString := n3.ReturnVCString()

	oString1 := n1.ReturnVCString()
	oString2 := n2.ReturnVCString()

	if an != 2 || bn != 3 || cn != 1 {
		t.Fatalf("Merge not as expected = %s , old = %s, %s", nString, oString1, oString2)
	} else if !n1.Compare(n3, Descendant) {
		failComparison(t, "Clocks not defined as Descendant: n1 = %s | n2 = %s", n1, n3)
	} else if !n2.Compare(n3, Descendant) {
		failComparison(t, "Clocks not defined as Descendant: n1 = %s | n2 = %s", n2, n3)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as concurrent: n1 = %s | n2 = %s", n1, n2)
	}
}

func TestCompareDiffLengthsNonConcurrent(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n2.Set("a", 1)
	n2.Set("b", 1)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks not defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks are defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks not defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks are defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareDiffLengthsConcurrent(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 2)
	n2.Set("a", 1)
	n2.Set("b", 1)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareIdenticalClocks(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n1.Set("b", 2)
	n1.Set("c", 3)
	n2.Set("a", 1)
	n2.Set("b", 2)
	n2.Set("c", 3)

	if !n1.Compare(n2, Equal) {
		failComparison(t, "Clocks not defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if !n2.Compare(n1, Equal) {
		failComparison(t, "Clocks not defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareSameLengthConcurrent(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n1.Set("b", 2)
	n1.Set("c", 3)
	n2.Set("a", 3)
	n2.Set("b", 2)
	n2.Set("c", 1)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareSameLengthNonConcurrent(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n1.Set("b", 2)
	n1.Set("c", 3)
	n2.Set("a", 2)
	n2.Set("b", 2)
	n2.Set("c", 3)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks not defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks are defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks not defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks are defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareNonIdenticalNames(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n1.Set("b", 2)
	n1.Set("c", 3)
	n2.Set("a", 1)
	n2.Set("b", 2)
	n2.Set("d", 3)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func TestCompareNonIdenticalNames2(t *testing.T) {
	n1 := New()
	n2 := New()

	n1.Set("a", 1)
	n1.Set("b", 1)

	n2.Set("b", 1)
	n2.Set("c", 1)
	n2.Set("d", 1)

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n1, n2)
	} else if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n1, n2)
	} else if !n1.Compare(n2, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n1, n2)
	}

	if n2.Compare(n1, Equal) {
		failComparison(t, "Clocks are defined as Equal: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Ancestor) {
		failComparison(t, "Clocks are defined as Ancestor: n1 = %s | n2 = %s", n2, n1)
	} else if n2.Compare(n1, Descendant) {
		failComparison(t, "Clocks are defined as Descendant: n1 = %s | n2 = %s", n2, n1)
	} else if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks not defined as Concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}

func failComparison(t *testing.T, failMessage string, clock1, clock2 VClock) {
	t.Fatalf(failMessage, clock1.ReturnVCString(), clock2.ReturnVCString())
}

func TestEncodeDecode(t *testing.T) {
	n := New()
	n.Set("a", 4)
	n.Set("b", 1)
	n.Set("c", 8)
	n.Set("d", 32)

	byteClock := n.Bytes()
	decoded, err := FromBytes(byteClock)

	if err != nil {
		t.Fatal(err)
	} else if !n.Compare(decoded, Equal) {
		nString := n.ReturnVCString()
		dString := decoded.ReturnVCString()
		t.Fatalf("decoded not the same as encoded enc = %s | dec = %s", nString, dString)
	}
}

func TestVCString(t *testing.T) {
	n := New()

	n.Set("a", 1)
	n.Set("b", 1)
	n.Set("c", 1)
	n.Set("d", 1)
	n.Set("e", 1)
	n.Set("f", 1)
	n.Set("g", 1)
	n.Set("h", 1)

	expected := "{\"a\":1, \"b\":1, \"c\":1, \"d\":1, \"e\":1, \"f\":1, \"g\":1, \"h\":1}"
	nString := n.ReturnVCString()

	if nString != expected {
		t.Fatalf("VC string %s not the same as expected %s", nString, expected)
	}
}

func Test2(t *testing.T) {
	n1 := New()
	n2 := New()

	// [map[node2:1 node3:1 node9:1]]
	n1.Set("A", 1)
	n1.Set("B", 1)
	n1.Set("C", 1)

	// map[node1:1 node3:4 node5:2 node6:1 node8:2 node9:3]
	n2.Set("B", 4)
	n2.Set("C", 3)
	n2.Set("D", 1)
	n2.Set("E", 1)

	fmt.Printf("n1 concurrent n2: %v\n", n1.Compare(n2, Concurrent))
	fmt.Printf("n2 concurrent n1: %v\n", n2.Compare(n1, Concurrent))
	fmt.Printf("n1 ancestor n2: %v\n", n1.Compare(n2, Ancestor))
	fmt.Printf("n2 ancestor n1: %v\n", n2.Compare(n1, Ancestor))
	fmt.Printf("n1 descendant n2: %v\n", n1.Compare(n2, Descendant))
	fmt.Printf("n2 descendant n1: %v\n", n2.Compare(n1, Descendant))
	fmt.Printf("n1 equal n2: %v\n", n1.Compare(n2, Equal))
	fmt.Printf("n2 equal n1: %v\n", n2.Compare(n1, Equal))

	if n1.Compare(n2, Descendant) {
		failComparison(t, "Clocks are defined as descendant: n1 = %s | n2 = %s", n1, n2)
	}

	if n1.Compare(n2, Ancestor) {
		failComparison(t, "Clocks are defined as ancestor: n1 = %s | n2 = %s", n2, n1)
	}

	if n1.Compare(n2, Equal) {
		failComparison(t, "Clocks are defined as equal: n1 = %s | n2 = %s", n1, n2)
	}

	if !n2.Compare(n1, Concurrent) {
		failComparison(t, "Clocks are not defined as concurrent: n1 = %s | n2 = %s", n2, n1)
	}
}
