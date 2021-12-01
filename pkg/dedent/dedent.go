package dedent

import (
	"strings"
)

// Indent inserts prefix in front of each line in s.
// Lines in s are delimited with \n.
//
func Indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// Dedent removes common leading space/tab strings from each line in s.
// Lines in s are delimited with \n.
// "common leading space/tab" means the exact same combination of spaces
// and tabs.
// The first line is special.
// If it does not start with space or tab, it will not be included in the
// calculation.
//
func Dedent(s string) string {
	lines := strings.Split(s, "\n")

	// find minimum leading spaces
	common := ""
	commonValid := false
	for i, line := range lines {
		if len(line) == 0 {
			continue // skip empty lines
		}

		if i == 0 {
			if !IsSpace(line[0]) {
				continue // skip first line if it doesn't start with space
			}
		}

		leading := LeadingWhite(line)

		if commonValid {
			common = CommonPrefix(common, leading)
		} else {
			common = leading
			commonValid = true
		}
	}

	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, common)
	}

	return strings.Join(lines, "\n")
}

// IsSpace is true if c is a space character or a tab.
//
func IsSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// LeadingWhite returns the prefix of s which consists of only spaces and tabs.
//
func LeadingWhite(s string) string {
	for i := 0; i < len(s); i++ {
		if !IsSpace(s[i]) {
			return s[0:i]
		}
	}
	return s
}

// CommonPrefix returns a string containing the common prefix of s1 and s2.
// Returns an empty string if there is no common prefix.
//
func CommonPrefix(s1, s2 string) string {
	shortest := len(s1)
	if len(s2) < shortest {
		shortest = len(s2)
	}

	var i int
	for i = 0; i < shortest; i++ {
		if s1[i] != s2[i] {
			break
		}
	}

	return s1[0:i]
}
