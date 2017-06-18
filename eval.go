package main

import "github.com/xtaci/goeval"

func EvalWithVals(vals map[string]interface{}) *goeval.Scope {
	s := goeval.NewScope()
	for key, value := range vals {
		s.Set(key, value)
	}
	return s
}
