package main

import "path"

func main() {
	raw := "zcong1993/test"
	gc, err := normalizeUrl(raw)
	exitWithError(err)
	zipPath, err := downloadGithubZipFile(*gc, "")
	exitWithError(err)
	err1 := extractFiles(zipPath, "", path.Join(TemplateHome, gc.Repo.Owner, gc.Repo.Name))
	checkErr(err1)
}
