## mergermarket/queryparser

[![GoDoc](https://godoc.org/github.com/mergermarket/queryparser?status.svg)](https://godoc.org/github.com/mergermarket/queryparser) [![Build Status](https://travis-ci.org/mergermarket/queryparser.svg?branch=master)](https://travis-ci.org/mergermarket/queryparser)


Package queryparser parses a Google-like search string to return a Query object.

The resulting Query object can be used as an interim representation between user input and something like ElasticSearch's Query DSL.

## Usage

A simple example:

```go
package main
import (
    "fmt"
    "github.com/mergermarket/queryparser"
    "strings"
)

func main() {
	query, _ := queryparser.NewParser(strings.NewReader(`(apple OR pear) AND (pie OR crumble)`)).Parse()
	fmt.Println(query)
}
```

output:
```bash
{Occur:MUST, SubQueries:[{Occur:SHOULD, Terms:["apple","pear"]},{Occur:SHOULD, Terms:["pie","crumble"]}]}

```
## What is Query?

Query is a simple but hopefully useful model of the intention behind the parsed string.

1. **Occur**: how to determine whether a document matches the query.
    1. MUST = all query components must match.
    2. SHOULD = at least one query component must match.
2. **Terms**: zero or more words or phrases
3. **SubQueries**: zero or more Query objects. So yeah wow like recursive, man.

## Syntax

The parser can handle arbitrary combinations of:

*  single word
*  phrase surrounded by quotes
*  any of the above combined with boolean operators "OR" and "AND"
*  any of the above grouped with parentheses

| Query component | Example | Intention |
| --- | --- | --- |
| single word | kayak | Match documents containing the word `kayak`|
| multiple words | kayak hammer | Match documents containing both `kayak` and `hammer`, in any order, not necessarily together. |
| " " | "roof slate" | Match documents containing the exact phrase `roof slate`. |
| ( ) | (baboon OR warthog) snake | Match documents containing at least one of `baboon` or `warthog`, and also `snake`. |
| AND | book AND pdf | Match both query components. Equivalent to `book pdf` |
| OR | tree OR shrub | Match either query component. |

Look at the [parser tests](./parser_test.go) for more examples.

## Operator precedence

Operators have the following precedence, in descending order:

| Operator |
|--------- |
| ( ) |
| AND |
| OR  |

For example, because AND has higher precedence than OR, the following queries:

```apple AND pear OR rhubarb```

```apple OR pear AND rhubarb```

evaluate as if written as follows:

```
(apple AND pear) OR rhubarb
```
```
apple OR (pear AND rhubarb)
```

## TODO

Not yet supported:

* "foo NOT bar" match only documents that contain "foo" but not "bar"
*  "-bar" match only documents that do not contain "bar"
