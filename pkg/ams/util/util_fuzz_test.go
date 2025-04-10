package util

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzQualifiedName(f *testing.F) {
	// we allow any printable utf8 character
	// and should include some special characters
	f.Add("a")
	f.Add("a....b...c")
	f.Add("auhdsiufa")
	f.Add("你好_jadpö世界\"")
	f.Add("ji...لسلام عليكم dsfüka&%$ ::.sdfaf?ßsd")
	f.Add("sadasd.asok\"a.b")
	f.Add("sad..\"...//  \\ssas")

	f.Fuzz(func(t *testing.T, input string) {

		// Split the input string into parts
		parts := randomSplit(input)

		for _, part := range parts {
			if part == "" {
				t.Errorf("part %q is empty", part)
			}
			if !utf8.ValidString(part) {
				t.Skip()
			}
		}
		if strings.Join(parts, "") != input {
			t.Skip()
			t.Errorf("randomSplit(%q) = %q; expected %q", input, strings.Join(parts, ""), input)
		}

		// Join the parts back together with a dot
		output := StringifyQualifiedName(parts)

		// Parse the output string back into parts
		parsedParts := ParseQualifiedName(output)

		// Check if the parsed parts are equal to the original parts
		if !reflect.DeepEqual(parts, parsedParts) {
			t.Errorf("Original   : %v", parts)
			t.Errorf("Stringified: %s", output)
			t.Errorf("Parsed     : %v", parsedParts)
		}
	})

}

func randomSplit(s string) []string {
	result := []string{}
	part := ""
	for _, c := range s {
		part += string(c)
		if rand.Intn(3) == 0 {
			result = append(result, part)
			part = ""
		}
	}
	if part != "" {
		result = append(result, part)
	}
	return result
}
