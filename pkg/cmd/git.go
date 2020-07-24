package cmd

import (
	"fmt"
	"log"
	"strings"
)

// Git keeps all dependency interface
type Git struct {
	PreviousTagGetter
	CommitGetter
}

// Commit stores basic data about git commit
type Commit struct {
	Hash    string
	Message string
}

type PreviousTagGetter func(name string, arg ...string) (string, error)
type CommitGetter func(name string, arg ...string) (string, error)

// New creates Git with default dependencies
func New() Git {
	return Git{
		PreviousTagGetter: Exec,
		CommitGetter:      Exec,
	}
}

// GetCommits gets all commits between current and previous tag
func (c Git) GetCommits(currentTag, previousTag, gitPath string) ([]Commit, error) {
	var commits []Commit
	r := fmt.Sprintf("%s...%s", currentTag, previousTag)
	log.Printf("[GIT] found tags: %s", r)

	out, err := c.CommitGetter("git", "-C", gitPath, "log", "--pretty=format:\"%H;%s\"", "--no-notes", r)
	if err != nil {
		return nil, err
	}

	resultLines := strings.Split(out, "\n")
	for _, line := range resultLines {
		l := strings.Split(line, ";")
		if len(l) > 1 {
			commits = append(commits, Commit{
				Hash:    l[0],
				Message: l[1],
			})
		}
	}

	return commits, nil
}

// GetPreviousTag tries to get one tag before given tag
func (c Git) GetPreviousTag(tag, gitPath string) (string, error) {
	out, err := c.PreviousTagGetter("git", "-C", gitPath, "describe", "--tags", "--abbrev=0", tag+"^")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil

}
