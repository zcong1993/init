package main

// fork from https://github.com/gruntwork-io/fetch/blob/master/github.go

// GitHubRepo is struct of github repo
type GitHubRepo struct {
	URL   string // The URL of the GitHub repo
	Owner string // The GitHub account name under which the repo exists
	Name  string // The GitHub repo name
	Token string // The personal access token to access this repo (if it's a private repo)
}

// GitHubCommit is struct of github commit
// Represents a specific git commit.
// Note that code using GitHub Commit should respect the following hierarchy:
// - CommitSha > BranchName > GitTag
// - Example: GitTag and BranchName are both specified; use the GitTag
// - Example: GitTag and CommitSha are both specified; use the CommitSha
// - Example: BranchName alone is specified; use BranchName
type GitHubCommit struct {
	Repo       GitHubRepo // The GitHub repo where this release lives
	GitTag     string     // The specific git tag for this release
	BranchName string     // If specified, indicates that this commit should be the latest commit on the given branch
	CommitSha  string     // If specified, indicates that this commit should be exactly this Git Commit SHA.
}
