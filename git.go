package git

import (
	"errors"
	"os/exec"
	"strings"
)

// Repo represents a git repo
type Repo struct {
	Path string
}

// PullResult is the result of a fetch & merge.
type PullResult struct {
	Repo    string
	Commits []string
}

// New returns a new repo-instance
func New(path string) (Repo, error) {

	// testing if path is a valid git +epository
	if RevParse(path, "@") == "" {
		err := errors.New(path + " is not a valid git repo")
		return Repo{}, err
	}

	return Repo{Path: path}, nil
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
func (r *Repo) Merge() error {
	cmd := exec.Command("git", "-C", r.Path, "merge", "--quiet")
	return cmd.Run() // returns error obj
}

// Fetch fetches changes from origin
func (r *Repo) Fetch() (string, error) {
	cmd := exec.Command("git", "-C", r.Path, "fetch")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Log returns unmerged commits
func (r *Repo) Log() (string, error) {
	cmd := exec.Command("git", "-C", r.Path, "log", "..origin", "--pretty=format:%h - %s (%an, %ar)")
	out, err := cmd.Output()
	return string(out), err
}

// Pull fetches and merges changes
func (r *Repo) Pull() PullResult {

	res := PullResult{Repo: r.Path}

	// fetch remote
	r.Fetch()

	// compare local and remote to check for changes
	local := RevParse(r.Path, "@")
	remote := RevParse(r.Path, "origin/master")

	// are there any changes?
	if local != remote {

		// save new commit-messages in result object
		commits, _ := r.Log()
		for _, c := range strings.Split(commits, "\n") {
			res.Commits = append(res.Commits, c)
		}

		// merge them
		r.Merge()

	}

	return res
}
