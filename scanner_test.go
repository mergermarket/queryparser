package queryparser

import (
	"strings"
	"testing"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok token
		lit string
	}{
		// Special tokens (eof, ws)
		{s: ``, tok: eof},
		{s: ` `, tok: ws, lit: " "},
		{s: "\t", tok: ws, lit: "\t"},
		{s: "\n", tok: ws, lit: "\n"},

		// Misc characters
		{s: `"`, tok: dquote, lit: `"`},

		// Identifiers
		{s: `foo`, tok: literal, lit: `foo`},
		{s: `Zx12_3U_-`, tok: literal, lit: `Zx12_3U_`},
		{s: `or`, tok: literal, lit: "or"},

		// Keywords
		{s: `OR`, tok: or, lit: "OR"},
		{s: `AND`, tok: and, lit: "AND"},
	}

	for i, tt := range tests {
		s := newScanner(strings.NewReader(tt.s))
		tok, lit := s.scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
