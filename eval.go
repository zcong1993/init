package main

import "github.com/xtaci/goeval"

// EvalWithVals return a sandbox with values
func EvalWithVals(vals map[string]interface{}) *goeval.Scope {
	s := goeval.NewScope()
	for key, value := range vals {
		s.Set(key, value)
	}
	return s
}
