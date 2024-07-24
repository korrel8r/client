// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

// package names parses and constructs query and class name strings.
package types

import (
	"fmt"
	"strings"
)

const sep = ":"

type Class struct{ Domain, Name string }

func ParseClass(class string) Class {
	c := Class{}
	s := strings.SplitN(class, sep, 2)
	if len(s) > 0 {
		c.Domain = s[0]
	}
	if len(s) > 1 {
		c.Name = s[1]
	}
	return c
}

func (c Class) String() string { return fmt.Sprintf("%v:%v", c.Domain, c.Name) }

type Query struct {
	Class Class
	Data  string
}

func ParseQuery(query string) Query {
	s := strings.SplitN(query, sep, 3)
	q := Query{}
	if len(s) > 0 {
		q.Class.Domain = s[0]
	}
	if len(s) > 1 {
		q.Class.Name = s[1]
	}
	if len(s) > 2 {
		q.Data = s[2]
	}
	return q
}

func (q Query) String() string { return fmt.Sprintf("%v:%v", q.Class, q.Data) }
