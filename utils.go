package main

import (
	"regexp"
	"fmt"
	"os"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

var TemplateHome string

func init() {
	home, err := homedir.Dir()
	checkErr(err)
	TemplateHome = filepath.Join(home, ".init-templates")
}


func normalizeUrl(url string) (*GitHubCommit, *InitError) {
	reg := regexp.MustCompile(`([^/]+)/([^#]+)(#(.+))?`)
	matches := reg.FindSubmatch([]byte(url))
	if len(matches) < 3 {
		return nil, newError(100, "invalid repo link")
	}
	gitCommit := &GitHubCommit{
		Repo:GitHubRepo{
			Owner:string(matches[1]),
			Name:string(matches[2]),
		},
		BranchName:"master",
	}
	if string(matches[4]) != "" {
		gitCommit.BranchName = string(matches[4])
	}
	return gitCommit, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
