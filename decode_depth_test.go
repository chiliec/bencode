package bencode

import (
	"bytes"
	"strings"
	"testing"
)

func nestedLists(n int) []byte {
	return []byte(strings.Repeat("l", n) + strings.Repeat("e", n))
}

func nestedDicts(n int) []byte {
	return []byte(strings.Repeat("d1:a", n) + "i0e" + strings.Repeat("e", n))
}

func TestDecodeModestNestingSucceeds(t *testing.T) {
	for _, n := range []int{1, 10, 50, DefaultMaxDepth} {
		var out interface{}
		if err := NewDecoder(bytes.NewReader(nestedLists(n))).Decode(&out); err != nil {
			t.Fatalf("depth %d: unexpected error: %v", n, err)
		}
	}
}

// Deeply nested input must return an error instead of overflowing the stack.
func TestDecodeDeepNestingReturnsError(t *testing.T) {
	for name, b := range map[string][]byte{
		"lists": nestedLists(1_000_000),
		"dicts": nestedDicts(1_000_000),
		"limit": nestedLists(DefaultMaxDepth + 1),
	} {
		var out interface{}
		err := NewDecoder(bytes.NewReader(b)).Decode(&out)
		if err == nil || !strings.Contains(err.Error(), "max nesting depth") {
			t.Fatalf("%s: expected a nesting depth error, got %v", name, err)
		}
	}
}

func TestSetMaxDepth(t *testing.T) {
	d := NewDecoder(bytes.NewReader(nestedLists(10)))
	d.SetMaxDepth(5)
	var out interface{}
	if err := d.Decode(&out); err == nil {
		t.Fatal("expected an error with a max depth of 5")
	}

	d = NewDecoder(bytes.NewReader(nestedLists(300)))
	d.SetMaxDepth(0)
	if err := d.Decode(&out); err != nil {
		t.Fatalf("a max depth of 0 should disable the limit: %v", err)
	}
}
