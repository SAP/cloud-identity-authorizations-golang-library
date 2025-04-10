package util

import (
	"reflect"
	"testing"
)

func TestStringifyReference(t *testing.T) {
	t.Run("No escaping needed", func(t *testing.T) {
		ref := []string{"test", "test2"}
		expected := "test.test2"
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}
	})
	t.Run("Contains double quote", func(t *testing.T) {
		ref := []string{"test", "test2\""}
		expected := "test.\"test2\\\"\""
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}

	})

	t.Run("Contains dot", func(t *testing.T) {
		ref := []string{"test", "test2."}
		expected := "test.\"test2.\""
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}
	})

	t.Run("Contains backslash", func(t *testing.T) {
		ref := []string{"test", "test2\\"}
		expected := "test.test2\\"
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}
	})

	t.Run("Contains chinese characters", func(t *testing.T) {
		ref := []string{"你好"}
		expected := "你好"
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}
	})

	t.Run("Contains part with only dots", func(t *testing.T) {
		ref := []string{".", "."}
		expected := "\".\".\".\""
		if got := StringifyQualifiedName(ref); got != expected {
			t.Errorf("StringifyReference() = %v, expected %v", got, expected)
		}
		parsed := ParseQualifiedName(expected)
		if !reflect.DeepEqual(ref, parsed) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", parsed, ref)
		}
	})

}

func TestParseQualifiedName(t *testing.T) {
	t.Run("No escaping needed", func(t *testing.T) {
		qualifiedName := "test.test2"
		expected := []string{"test", "test2"}
		if got := ParseQualifiedName(qualifiedName); !reflect.DeepEqual(got, expected) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", got, expected)
		}
	})
	t.Run("Contains double quote", func(t *testing.T) {
		qualifiedName := "$app.\"\\\"quoted2\\\"\".findme"
		expected := []string{"$app", "\"quoted2\"", "findme"}
		if got := ParseQualifiedName(qualifiedName); !reflect.DeepEqual(got, expected) {
			t.Errorf("ParseQualifiedName() = %v, expected %v", got, expected)
		}
	})

}
