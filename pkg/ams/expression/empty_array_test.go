package expression

import (
	"testing"
)

func TestEmptyArray(t *testing.T) {
	empty := EmptyArray{}
	if !empty.IsEmpty() {
		t.Errorf("Expected empty array to be empty")
	}
	if empty.Contains(String("a")) {
		t.Errorf("Expected empty array to not contain any elements")
	}
	if len(empty.Elements()) != 0 {
		t.Errorf("Expected empty array to have no elements")
	}
	want := "[]"
	got := empty.String()
	if got != want {
		t.Errorf("Expected empty array string representation to be %s, got %s", want, got)
	}

	if empty.equals(StringArray{"a"}) {
		t.Errorf("Expected empty array to not equal non-empty array")
	}

	if empty.LessThan(empty) {
		t.Errorf("Expected empty array to not be less than itself")
	}
	if empty.Evaluate(nil) != empty {
		t.Errorf("Expected empty array to evaluate to itself")
	}
}
