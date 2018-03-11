package bloomflt

import (
	"fmt"
	"math"
	"testing"
)

func TestNewMK(t *testing.T) {
	b := NewMK(64, 2)
	if b == nil {
		t.Errorf("NewMK(64, 2) = nil, expected a struct")
	}
}

func TestCalcOptimalMK(t *testing.T) {
	gotM, gotK := CalcOptimalMK(216553, 0.01)
	wantM := 2075673
	if gotM != wantM {
		t.Errorf("CalcOptimalMK(216553, 0.01) = %v, want %v", gotM, wantM)
	}
	wantK := 7
	if gotK != wantK {
		t.Errorf("CalcOptimalMK(216553, 0.01) = %v, want %v", gotK, wantK)
	}
}

func TestHighFalsePositiveRate(t *testing.T) {
	b := New(100, 0.01)

	ok := b.ContainsString("SomeValue")
	if ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", "SomeValue", ok, false)
	}

	b.AddString("SomeValue")
	b.AddString("AnotherValue")

	ok = b.ContainsString("SomeValue")
	if !ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", "SomeValue", ok, true)
	}

	ok = b.ContainsString("AnotherValue")
	if !ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", "AnotherValue", ok, true)
	}
}

func TestLowFalsePositiveRate(t *testing.T) {
	b := New(1000, 0.00001)

	value := "AnotherValue"
	ok := b.ContainsString(value)
	if ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, false)
	}

	b.AddString(value)
	ok = b.ContainsString(value)
	if !ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, true)
	}
}

func TestMinimal(t *testing.T) {
	b := New(0, 0.01)

	value := "AnotherValue"
	ok := b.ContainsString(value)
	if ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, false)
	}

	b.AddString(value)
	ok = b.ContainsString(value)
	if !ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, true)
	}
}

func TestMaximum(t *testing.T) {
	b := New(math.MaxInt32, 0.001)

	value := "AnotherValue"
	ok := b.ContainsString(value)
	if ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, false)
	}

	b.AddString(value)
	ok = b.ContainsString(value)
	if !ok {
		t.Errorf("b.ContainsString(%q) = %v, want %v", value, ok, true)
	}
}

func TestUInt32(t *testing.T) {
	b := New(100, 0.01)

	value := uint32(math.MaxUint32 / 2)

	ok := b.ContainsUInt32(value)
	if ok {
		t.Errorf("b.ContainsString(%v) = %v, want %v", value, ok, false)
	}

	b.AddUInt32(value)
	ok = b.ContainsUInt32(value)
	if !ok {
		t.Errorf("b.ContainsString(%v) = %v, want %v", value, ok, true)
	}
}

func TestUInt64(t *testing.T) {
	b := New(100, 0.01)

	value := uint64(math.MaxUint64 / 2)

	ok := b.ContainsUInt64(value)
	if ok {
		t.Errorf("b.ContainsString(%v) = %v, want %v", value, ok, false)
	}

	b.AddUInt64(value)
	ok = b.ContainsUInt64(value)
	if !ok {
		t.Errorf("b.ContainsString(%v) = %v, want %v", value, ok, true)
	}
}

func Example() {
	// We expect the set to contains up to 100 elements with acceptable false-positive rate of 0.01%
	b := New(100, 0.01)

	b.AddString("value1")
	if b.ContainsString("value1") {
		fmt.Println("The set now has 'value1'.")
	}

	someID := uint64(123)
	b.AddUInt64(someID)
	if b.ContainsUInt64(someID) {
		fmt.Printf("The set now has ID %v.", someID)
	}

	// Output: The set now has 'value1'.
	// The set now has ID 123.
}
