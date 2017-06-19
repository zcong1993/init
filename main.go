package main

import (
	"os"
	"path/filepath"
)

func main() {
	cli := &CLI{os.Stdout, os.Stderr}
	os.Exit(cli.Run(os.Args))
}

func bak() {
	url := "zcong1993/test"
	gc, err := normalizeUrl(url)
	exitWithError(err)
	zipPath, err := downloadGithubZipFile(*gc, "")
	exitWithError(err)
	err1 := extractFiles(zipPath, "", filepath.Join(TemplateHome, gc.Repo.Owner, gc.Repo.Name))
	checkErr(err1)
}
