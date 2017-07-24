package queryparser

import (
	"fmt"

	"strings"
)

// Query represents a query with one or more terms, how to decide matching, and zero or more subqueries
type Query struct {
	Terms      []string
	Occur      Occur
	SubQueries []Query
}

// Term represents a word or phrase
type Term string

// Occur represents whether we want documents containing all terms and matching all subqueries (MUST) or at least one term or matching at least one subquery (SHOULD)
type Occur int

const (
	// MUST means we want documents containing all terms and matching all subqueries
	MUST Occur = iota
	// SHOULD means we want documents containing at least one term or matching at least one subquery
	SHOULD
)

func (occur Occur) string() string {
	switch occur {
	case MUST:
		return "MUST"
	case SHOULD:
		return "SHOULD"
	default:
		return ""
	}

}

func termsString(ts []string) string {
	if len(ts) == 0 {
		return ""
	}
	sts := make([]string, len(ts))
	for i, t := range ts {
		sts[i] = fmt.Sprintf(`"%s"`, t)
	}
	return fmt.Sprintf(", Terms:[%s]", strings.Join(sts, ","))
}

func subqueriesString(qs []Query) string {
	if len(qs) == 0 {
		return ""
	}
	sqs := make([]string, len(qs))
	for i, t := range qs {
		sqs[i] = fmt.Sprintf(`%s`, t)
	}
	return fmt.Sprintf(", SubQueries:[%s]", strings.Join(sqs, ","))
}

func (q Query) String() string {
	return fmt.Sprintf("{Occur:%s%s%s}", q.Occur.string(), termsString(q.Terms), subqueriesString(q.SubQueries))
}
