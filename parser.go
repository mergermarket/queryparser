package queryparser

import (
	"fmt"
	"io"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	s   *scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newScanner(r)}
}

// Parse parses a query string.
func (p *Parser) Parse() (*Query, error) {
	var phraseWords []string
	var insideAPhrase bool
	var currentSubquery *Query
	var currentSubqueryParent *Query
	query := &Query{Occur: MUST}
	var currentTerms = &query.Terms
	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case dquote:
			// If this quote is end of a phrase
			if insideAPhrase {
				*currentTerms = append(*currentTerms, strings.Join(phraseWords, " "))
				phraseWords = nil
				insideAPhrase = false
				// Otherwise assume it is start of a phrase
			} else {
				insideAPhrase = true
			}
		case lparen:
			// Ignore parentheses inside phrases
			if !insideAPhrase {
				if currentSubquery == nil {
					query.SubQueries = append(query.SubQueries, Query{Occur: MUST})
					currentSubqueryParent = query
					currentSubquery = &query.SubQueries[0]

				} else {
					currentSubquery.SubQueries = append(currentSubquery.SubQueries, Query{Occur: MUST})
					currentSubqueryParent = currentSubquery
					currentSubquery = &currentSubquery.SubQueries[len(currentSubquery.SubQueries)-1]
				}
				currentTerms = &currentSubquery.Terms
			}
		case rparen:
			// Ignore parentheses inside phrases
			if !insideAPhrase {
				currentSubquery = currentSubqueryParent
				currentTerms = &currentSubquery.Terms
			}
		case literal:
			if insideAPhrase {
				phraseWords = append(phraseWords, lit)
			} else {
				*currentTerms = append(*currentTerms, lit)
			}
		case and:
			if insideAPhrase {
				phraseWords = append(phraseWords, lit)
			} else {
				if len(*currentTerms) > 0 {
					if currentSubquery != nil {
						currentSubquery.Occur = MUST
					} else {
						query.Occur = MUST
					}
				}
			}
		case or:
			if insideAPhrase {
				phraseWords = append(phraseWords, lit)
			} else {
				if currentSubquery != nil {
					currentSubquery.Occur = SHOULD
					// Because and has higher precedence than or, implicitly create subquery
				} else if len(*currentTerms) > 1 && query.Occur == MUST {
					query.SubQueries = append(query.SubQueries, Query{Occur: SHOULD})
					currentSubqueryParent = query
					currentSubquery = &query.SubQueries[0]
					currentTerms = &currentSubquery.Terms
				} else {
					query.Occur = SHOULD
				}
			}
		case eof:
			return query, nil
		default:
			fmt.Printf("default: tok=%v lit=%v\n", tok, lit)
		}
	}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok token, lit string) {
	tok, lit = p.scan()
	if tok == ws {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
