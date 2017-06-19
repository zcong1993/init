package main

import (
	"os"
	"path/filepath"
)

func main() {
	cli := &CLI{os.Stdout, os.Stderr}
	os.Exit(cli.Run(os.Args))
	//bak()
}

func bak() {
	data, _ := GetConfig("./config.json")
	CopyFile(filepath.Join(TemplateHome, "zcong1993", "test", "template", "README.md"), filepath.Join(".", "test", "fx.md"), data)
}
