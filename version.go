package release

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Versioner interface {
	Version() (string, error)
}

type providedVersion struct {
	version string
}

// NewProvidedVersion creates a versioner that return the
// given version
func NewProvidedVersion(major, minor, patch int) Versioner {
	return &providedVersion{
		version: fmt.Sprintf("v%d.%d.%d", major, minor, patch),
	}
}

// Version returns the provided version
func (pv *providedVersion) Version() (string, error) {
	return pv.version, nil
}

type gitVersion struct {
	repositoryPath string
	branch         string
	major          int
}

// NewGitHistoryVersion creates a versioner that can extract
// a version from a git commit history
func NewGitHistoryVersion(repositoryPath string, branch string, major int) Versioner {
	return &gitVersion{
		repositoryPath: repositoryPath,
		branch:         branch,
		major:          major,
	}
}

// Version returns the generated version using the
// git history
func (gv *gitVersion) Version() (string, error) {
	r, err := git.PlainOpen(gv.repositoryPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to open git repository")
	}
	w, err := r.Worktree()
	if err != nil {
		return "", errors.Wrap(err, "failed to load git work tree")
	}
	err = w.Checkout(
		&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", gv.branch)),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to checkout branch: %s", gv.branch))
	}

	iter, err := r.Log(&git.LogOptions{})
	defer iter.Close()
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch the git log")
	}

	var minor, patch int
	counterFn := func(commit *object.Commit) error {
		parentCount := 0
		parents := commit.Parents()
		defer parents.Close()
		_ = parents.ForEach(func(*object.Commit) error {
			parentCount++
			return nil
		})
		if parentCount == 1 {
			patch++
		} else if parentCount == 2 { // A commit with two parents is a merge commit
			minor++
		}
		return nil
	}
	_ = iter.ForEach(counterFn)
	return fmt.Sprintf("v%d.%d.%d", gv.major, minor, patch), nil
}
