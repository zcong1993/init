package main

import (
	"github.com/buger/jsonparser"
	"fmt"
)

func parse(data []byte) {
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
		return nil
	}, "prompts")
}
