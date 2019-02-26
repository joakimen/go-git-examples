package git_test

import (
	"testing"

	"github.com/krystah/git"
)

func TestIsValidRepo(t *testing.T) {
	isValid := git.IsValidRepo(".")

	if !isValid {
		t.Error(". is not a valid repository")
	}
}
