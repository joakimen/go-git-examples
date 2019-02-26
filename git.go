package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// IsValidRepo returns a new repo-instance
func IsValidRepo(path string) bool {

	// testing if path is a valid git +epository
	if RevParse(path, "@") == "" {
		return false
	}

	return true
}

// RevParse attempts to get the revision hash of the specified path
func RevParse(path, rev string) string {
	cmd := exec.Command("git", "-C", path, "rev-parse", rev)

	// an error here simply means a revision couldn't be parsed from the
	// specified directory. this is fine
	out, _ := cmd.Output()

	return string(out)
}

// Merge merges changes from origin
func Merge(repo string) error {
	cmd := exec.Command("git", "-C", repo, "merge", "--quiet")
	out, err := cmd.Output()

	fmt.Println(string(out))

	if err != nil {
		fmt.Println("merge failed")
		return err
	}

	fmt.Println("succeeded")
	return nil
}

// Fetch fetches changes from origin
func Fetch(repo string) error {
	cmd := exec.Command("git", "-C", repo, "fetch")
	err := cmd.Run()
	return err
}

// Log returns unmerged commits
func Log(repo string, revisionRange string) (commits []string, err error) {
	cmd := exec.Command("git", "-C", repo, "log", revisionRange, "--pretty=format:%h - %s (%an, %ar)")
	out, err := cmd.Output()
	if err != nil {
		return commits, err
	}

	// populate commits from log
	for _, c := range strings.Split(string(out), "\n") {
		commits = append(commits, c)
	}
	return commits, err
}

// Pull fetches and merges changes
func Pull(repo string) ([]string, error) {

	var commits []string

	// fetch remote
	err := Fetch(repo)
	if err != nil {
		return commits, errors.Wrap(err, "git fetch")
	}

	// compare local and remote to check for changes
	local := RevParse(repo, "@")
	remote := RevParse(repo, "origin/master")

	// are there any changes?
	if local != remote {

		// save new commit-messages in result object
		commits, err = Log(repo, "@..origin")
		if err != nil {
			return commits, errors.Wrap(err, "git log")
		}

		// merge them
		Merge(repo)

	}

	return commits, nil
}
