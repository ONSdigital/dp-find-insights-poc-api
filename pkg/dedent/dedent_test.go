package dedent_test

import (
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/dedent"
)

func Test_IsSpace(t *testing.T) {
	var tests = []struct {
		c    byte
		want bool
	}{
		{'a', false},
		{' ', true},
		{'\t', true},
	}

	for _, test := range tests {
		got := dedent.IsSpace(test.c)
		if got != test.want {
			t.Errorf("%c: %t, want %t\n", test.c, got, test.want)
		}
	}
}

func Test_LeadingWhite(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"", ""},
		{" ", " "},
		{"\t", "\t"},
		{" a", " "},
		{"\ta", "\t"},
		{" \t a", " \t "},
		{"a", ""},
	}

	for _, test := range tests {
		got := dedent.LeadingWhite(test.in)
		if got != test.want {
			t.Errorf("%q: %q, want %q\n", test.in, got, test.want)
		}
	}
}

func Test_CommonPrefix(t *testing.T) {
	var tests = []struct {
		s1   string
		s2   string
		want string
	}{
		{"", "", ""},
		{"a", "a", "a"},
		{"a", "ab", "a"},
		{"ab", "ac", "a"},
		{"x", "y", ""},
		{"xa", "y", ""},
		{"aba", "abcd", "ab"},
		{"abc", "abc", "abc"},
	}

	for _, test := range tests {
		got := dedent.CommonPrefix(test.s1, test.s2)
		if got != test.want {
			t.Errorf("%q,%q: %q, want %q\n", test.s1, test.s2, got, test.want)
		}
	}
}
func Test_Dedent(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{
			"",
			"",
		},
		{
			"\n",
			"\n",
		},
		// first line is special
		{
			`no leading spaces on first line
			 but some on next lines
			 `,
			"no leading spaces on first line\nbut some on next lines\n",
		},
		// blank lines not counted
		{
			`
			blank lines

			are not counted
			`,
			"\nblank lines\n\nare not counted\n",
		},
		// no newline on last "line"
		{
			" a\n b",
			"a\nb",
		},
	}

	for _, test := range tests {
		got := dedent.Dedent(test.in)
		if got != test.want {
			t.Errorf("input:\n%s\ngot:\n%s\nwant:\n%s\n", test.in, got, test.want)
		}
	}
}

func Test_Indent(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"", "    "},
		{"\n", "    \n    "},
		{"line", "    line"},
		{"line\n", "    line\n    "},
	}

	for _, test := range tests {
		got := dedent.Indent(test.in, "    ")
		if got != test.want {
			t.Errorf("%q: %q, want %q\n", test.in, got, test.want)
		}
	}
}
