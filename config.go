package main

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/tj/go-prompt"
	"io/ioutil"
)

// Config is type of init.json of template
type Config struct {
	Prompts map[string]Prompt
	Filters map[string]string
}

// Prompt is type of prompts in init.json
type Prompt struct {
	Message string
	Type    string
	Default string
}

// Cfg is a struct of some information of config
type Cfg struct {
	Path      string
	RawConfig []byte
	Config    *Config
}

// NewConfig create a Cfg from config file
func NewConfig(p string) (*Cfg, error) {
	rawConfig, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, err
	}
	return &Cfg{Path: p, RawConfig: rawConfig, Config: &config}, nil
}

// GetPrompts provide data by prompts
func (cfg *Cfg) GetPrompts() (map[string]interface{}, error) {
	config := cfg.Config
	res := map[string]interface{}{}
	for key, val := range config.Prompts {
		if val.Type == "confirm" {
			res[key] = prompt.Confirm(val.Message + " y/n ")
		}
		if val.Type == "string" {
			if val.Default == "" {
				res[key] = prompt.StringRequired(val.Message + " : ")
			} else {
				res[key] = prompt.String(val.Message+"(Default is %s) : ", val.Default)
				if res[key] == "" {
					res[key] = val.Default
				}
			}
		}
		if val.Type == "list" {
			list := []string{}
			jsonparser.ArrayEach(cfg.RawConfig, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				list = append(list, string(value))
			}, "prompts", key, "choices")
			i := prompt.Choose(val.Message, list)
			res[key] = list[i]
		}
	}
	return res, nil
}
