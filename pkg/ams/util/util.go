package util

import (
	"strings"
)

// reverse to StringifyQualifiedName
//
// if the provided string is not a valid qualified name, the return value is not reliable
func ParseQualifiedName(qualifiedName string) []string {
	result := make([]string, 0)

	currentPart := ""
	const (
		startofPart = iota
		quotedPart
		unquotedPart
		escapedRune
		endOfQuotedPart
	)
	state := startofPart
	for _, c := range qualifiedName {

		switch state {
		case startofPart:
			switch c {
			case '"':
				state = quotedPart
			default:
				currentPart = string(c)
				state = unquotedPart
			}
		case quotedPart:
			switch c {
			case '\\':
				state = escapedRune
			case '"':
				result = append(result, currentPart)
				currentPart = ""
				state = endOfQuotedPart
			default:
				currentPart += string(c)
			}
		case unquotedPart:
			switch c {
			case '.':
				result = append(result, currentPart)
				currentPart = ""
				state = startofPart
			default:
				currentPart += string(c)
			}
		case escapedRune:
			currentPart += string(c)
			state = quotedPart
		case endOfQuotedPart:
			state = startofPart

		}
	}
	if currentPart != "" {
		result = append(result, currentPart)
	}

	return result

}

// returns the string representation of a qualified name
//
// the parts are concatenated with a dot. Parts containing a dot or a quote will be surrounded by quotes and inner quotes and backslashes will be escaped with a backslash
func StringifyQualifiedName(ref []string) string {
	escaped := make([]string, len(ref))
	for i, r := range ref {
		if strings.ContainsAny(r, ".\"") {
			escaped[i] = escapePart(r)
		} else {
			escaped[i] = r
		}
	}
	return strings.Join(escaped, ".")
}

func escapePart(part string) string {

	result := strings.ReplaceAll(part, "\\", "\\\\")
	result = strings.ReplaceAll(result, "\"", "\\\"")
	return "\"" + result + "\""
}
