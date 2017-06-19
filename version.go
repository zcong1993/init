package main

import "fmt"

// Name is cli name
const Name = "init"

// Version is cli current version
const Version = "v0.0.1"

// GitCommit is cli current git commit hash
var GitCommit string

// ShowVersion is handler for version command
func ShowVersion() {
	version := fmt.Sprintf("%s version %s", Name, Version)
	if len(GitCommit) != 0 {
		version += fmt.Sprintf(" (%s)", GitCommit)
	}
	fmt.Println(version)
}
