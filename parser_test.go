package queryparser_test

import (
	"fmt"
	"github.com/mergermarket/queryparser"
	"reflect"
	"strings"
	"testing"
)

func ExampleParser_Parse() {
	query, _ := queryparser.NewParser(strings.NewReader(`(apple OR pear) AND (pie OR crumble)`)).Parse()
	fmt.Println(query)
	// Output:
	// {Occur:MUST, SubQueries:[{Occur:SHOULD, Terms:["apple","pear"]},{Occur:SHOULD, Terms:["pie","crumble"]}]}
}

// Ensure the parser can parse strings into Queries.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s     string
		query *queryparser.Query
		err   string
	}{
		{
			s: `foo`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo"},
			},
		},
		{
			s: `foo bar`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo", "bar"},
			},
		},
		{
			s: `"foo bar"`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo bar"},
			},
		},
		{
			s: `“curved quotes”`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"curved quotes"},
			},
		},
		{
			s: `foo OR bar`,
			query: &queryparser.Query{
				Occur: queryparser.SHOULD,
				Terms: []string{"foo", "bar"},
			},
		},
		{
			s: `foo or bar`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo", "or", "bar"},
			},
		},
		{
			s: `(foo AND boo) OR bar`,
			query: &queryparser.Query{
				Occur: queryparser.SHOULD,
				Terms: []string{"bar"},
				SubQueries: []queryparser.Query{{
					Occur: queryparser.MUST,
					Terms: []string{"foo", "boo"},
				}},
			},
		},
		{
			s: `foo boo OR bar`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo", "boo"},
				SubQueries: []queryparser.Query{{
					Occur: queryparser.SHOULD,
					Terms: []string{"bar"},
				}},
			},
		},
		{
			s: `foo AND (boo OR bar)`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo"},
				SubQueries: []queryparser.Query{{
					Occur: queryparser.SHOULD,
					Terms: []string{"boo", "bar"},
				}},
			},
		},
		{
			s: `foo (boo OR bar)`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"foo"},
				SubQueries: []queryparser.Query{{
					Occur: queryparser.SHOULD,
					Terms: []string{"boo", "bar"},
				}},
			},
		},
		{
			s: `(foo OR boo) AND bar`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"bar"},
				SubQueries: []queryparser.Query{{
					Occur: queryparser.SHOULD,
					Terms: []string{"foo", "boo"},
				}},
			},
		},
		{
			s: `(a AND b) OR (c AND d) OR (e AND f)`,
			query: &queryparser.Query{
				Occur: queryparser.SHOULD,
				SubQueries: []queryparser.Query{
					{
						Occur: queryparser.MUST,
						Terms: []string{"a", "b"},
					},
					{
						Occur: queryparser.MUST,
						Terms: []string{"c", "d"},
					},
					{
						Occur: queryparser.MUST,
						Terms: []string{"e", "f"},
					},
				},
			},
		},
		{
			s: `"parentheses (inside quotes) are ignored"`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				Terms: []string{"parentheses inside quotes are ignored"},
			},
		},
		{
			s: `("quotes inside" AND parentheses) OR ("are handled well")`,
			query: &queryparser.Query{
				Occur: queryparser.SHOULD,
				SubQueries: []queryparser.Query{
					{
						Occur: queryparser.MUST,
						Terms: []string{"quotes inside", "parentheses"},
					},
					{
						Occur: queryparser.MUST,
						Terms: []string{"are handled well"},
					},
				},
			},
		},
		{
			s: `("unbalanced quotes inside AND parentheses) OR ("are handled as well as possible")`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				SubQueries: []queryparser.Query{
					{
						Occur: queryparser.MUST,
						Terms: []string{"unbalanced quotes inside AND parentheses OR", "are", "handled", "as", "well", "as", "possible"},
					},
				},
			},
		},
		{
			s: `(((a`,
			query: &queryparser.Query{
				Occur: queryparser.MUST,
				SubQueries: []queryparser.Query{
					{
						Occur: queryparser.MUST,
						SubQueries: []queryparser.Query{
							{
								Occur: queryparser.MUST,
								SubQueries: []queryparser.Query{
									{
										Occur: queryparser.MUST,
										Terms: []string{"a"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		query, err := queryparser.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.query, query) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%v\n\ngot=%v\n\n", i, tt.s, tt.query, query)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
