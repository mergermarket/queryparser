package queryparser

// token represents a lexical token.
type token int

const (
	// Special tokens
	illegal token = iota
	eof
	ws

	// Literals
	literal // main

	// Misc characters
	lparen // (
	rparen // )
	dquote // "

	// Keywords
	or
	and
)
