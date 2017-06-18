package main

import (
	"encoding/json"
	"io/ioutil"
	"github.com/tj/go-prompt"
	"github.com/buger/jsonparser"
)

type Config struct {
	Prompts map[string]Prompt
}

type Prompt struct {
	Message string
	Type string
	Default string
}

func GetConfig(configPath string) (map[string]interface{}, error)  {
	rawConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	for key, val := range config.Prompts {
		if val.Type == "confirm" {
			res[key] = prompt.Confirm(val.Message + " y/n ")
		}
		if val.Type == "string" {
			if val.Default == "" {
				res[key] = prompt.StringRequired(val.Message + " : ")
			} else {
				res[key] = prompt.String(val.Message + "(Default is %s) : ", val.Default)
				if res[key] == "" {
					res[key] = val.Default
				}
			}
		}
		if val.Type == "list" {
			list := []string{}
			jsonparser.ArrayEach(rawConfig, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				list = append(list, string(value))
			}, "prompts", key, "choices")
			i := prompt.Choose(val.Message, list)
			res[key] = list[i]
		}
	}
	return res, nil
}
