package main

import (
	"regexp"
	"fmt"
	"os"
)

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
	fmt.Printf("%+v", gitCommit)
	return gitCommit, nil
}

func exitWithError(err *InitError) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(err.errorCode)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
