package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// fork from https://github.com/gruntwork-io/fetch/blob/master/file.go

// Download the zip file at the given URL to a temporary local directory.
// Returns the absolute path to the downloaded zip file.
// IMPORTANT: You must call "defer os.RemoveAll(dir)" in the calling function when done with the downloaded zip file!
func downloadGithubZipFile(gitHubCommit GitHubCommit, gitHubToken string) (string, *InitError) {

	var zipFilePath string

	// Create a temp directory
	// Note that ioutil.TempDir has a peculiar interface. We need not specify any meaningful values to achieve our
	// goal of getting a temporary directory.
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return zipFilePath, wrapError(err)
	}

	// Download the zip file, possibly using the GitHub oAuth Token
	httpClient := &http.Client{}
	req, err := MakeGitHubZipFileRequest(gitHubCommit, gitHubToken)
	if err != nil {
		return zipFilePath, wrapError(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return zipFilePath, wrapError(err)
	}
	if resp.StatusCode != 200 {
		return zipFilePath, newError(FAILED_TO_DOWNLOAD_FILE, fmt.Sprintf("Failed to download file at the url %s. Received HTTP Response %d.", req.URL.String(), resp.StatusCode))
	}
	if resp.Header.Get("Content-Type") != "application/zip" {
		return zipFilePath, newError(FAILED_TO_DOWNLOAD_FILE, fmt.Sprintf("Failed to download file at the url %s. Expected HTTP Response's \"Content-Type\" header to be \"application/zip\", but was \"%s\"", req.URL.String(), resp.Header.Get("Content-Type")))
	}

	// Copy the contents of the downloaded file to our empty file
	respBodyBuffer := new(bytes.Buffer)
	respBodyBuffer.ReadFrom(resp.Body)
	err = ioutil.WriteFile(filepath.Join(tempDir, "repo.zip"), respBodyBuffer.Bytes(), 0644)
	if err != nil {
		return zipFilePath, wrapError(err)
	}

	zipFilePath = filepath.Join(tempDir, "repo.zip")

	return zipFilePath, nil
}

// Decompress the file at zipFileAbsPath and move only those files under filesToExtractFromZipPath to localPath
func extractFiles(zipFilePath, filesToExtractFromZipPath, localPath string) error {

	// Open the zip file for reading.
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer func(p string) {
		os.RemoveAll(filepath.Dir(p))
	}(zipFilePath)
	defer r.Close()

	// pathPrefix represents the portion of the local file path we will ignore when copying the file to localPath
	// E.g. full path = fetch-test-public-0.0.3/folder/file1.txt
	//      path prefix = fetch-test-public-0.0.3
	//      file that will eventually get written = <localPath>/folder/file1.txt

	// By convention, the first file in the zip file is the top-level directory
	pathPrefix := r.File[0].Name

	// Add the path from which we will extract files to the path prefix so we can exclude the appropriate files
	pathPrefix = filepath.Join(pathPrefix, filesToExtractFromZipPath)

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {

		// If the given file is in the filesToExtractFromZipPath, proceed
		if strings.Index(f.Name, pathPrefix) == 0 {

			if f.FileInfo().IsDir() {
				// Create a directory
				os.MkdirAll(filepath.Join(localPath, strings.TrimPrefix(f.Name, pathPrefix)), 0777)
			} else {
				// Read the file into a byte array
				readCloser, err := f.Open()
				if err != nil {
					return fmt.Errorf("Failed to open file %s: %s", f.Name, err)
				}

				byteArray, err := ioutil.ReadAll(readCloser)
				if err != nil {
					return fmt.Errorf("Failed to read file %s: %s", f.Name, err)
				}

				// Write the file
				err = ioutil.WriteFile(filepath.Join(localPath, strings.TrimPrefix(f.Name, pathPrefix)), byteArray, 0644)
				if err != nil {
					return fmt.Errorf("Failed to write file: %s", err)
				}
			}
		}
	}

	return nil
}

// Return an HTTP request that will fetch the given GitHub repo's zip file for the given tag, possibly with the gitHubOAuthToken in the header
// Respects the GitHubCommit hierachy as defined in the code comments for GitHubCommit (e.g. GitTag > CommitSha)
func MakeGitHubZipFileRequest(gitHubCommit GitHubCommit, gitHubToken string) (*http.Request, error) {
	var request *http.Request

	// This represents either a commit, branch, or git tag
	var gitRef string
	if gitHubCommit.CommitSha != "" {
		gitRef = gitHubCommit.CommitSha
	} else if gitHubCommit.BranchName != "" {
		gitRef = gitHubCommit.BranchName
	} else if gitHubCommit.GitTag != "" {
		gitRef = gitHubCommit.GitTag
	} else {
		return request, fmt.Errorf("Neither a GitCommitSha nor a GitTag nor a BranchName were specified so impossible to identify a specific commit to download.")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", gitHubCommit.Repo.Owner, gitHubCommit.Repo.Name, gitRef)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return request, wrapError(err)
	}

	if gitHubToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", gitHubToken))
	}

	return request, nil
}

func DownloadAndExtract(gc *GitHubCommit) (string, *InitError) {
	dest := filepath.Join(TemplateHome, gc.Repo.Owner, gc.Repo.Name)
	if _, err := os.Stat(dest); err == nil {
		err = os.RemoveAll(dest)
		if err != nil {
			return "", wrapError(err)
		}
	}
	zipPath, err := downloadGithubZipFile(*gc, "")
	if err != nil {
		return "", err
	}
	err1 := extractFiles(zipPath, "", dest)
	if err1 != nil {
		return "", wrapError(err1)
	}
	return dest, nil
}
