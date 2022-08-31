package cmd

import (
	"fmt"
	"strings"

	pslog "github.com/psmarcin/jira-versioner/pkg/log"
)

const endOfCommit = "==EOC=="

// Git keeps all dependency interface
type Git struct {
	PreviousTagGetter
	CommitGetter

	log pslog.Logger
}

// Commit stores basic data about git commit
type Commit struct {
	Hash    string
	Message string
}

type PreviousTagGetter func(name string, arg ...string) (string, error)
type CommitGetter func(name string, arg ...string) (string, error)

// New creates Git with default dependencies
func New(log pslog.Logger) Git {
	return Git{
		PreviousTagGetter: Exec,
		CommitGetter:      Exec,
		log:               log,
	}
}

// GetCommits gets all commits between current and previous tag
func (c Git) GetCommits(currentTag, previousTag, gitPath string) ([]Commit, error) {
	var commits []Commit
	r := fmt.Sprintf("%s...%s", currentTag, previousTag)
	c.log.Infof("[GIT] found tags: %s", r)

	format := fmt.Sprintf(`--pretty=format:%%H;%%s %%b%s`, endOfCommit)
	out, err := c.CommitGetter("git", "-C", gitPath, "log", format, "--no-notes", r)
	if err != nil {
		return nil, err
	}

	resultLines := strings.Split(out, endOfCommit)
	for _, line := range resultLines {
		l := strings.Split(line, ";")
		if len(l) > 1 {
			commits = append(commits, Commit{
				Hash:    strings.ReplaceAll(l[0], "\n", ""),
				Message: strings.ReplaceAll(l[1], "\n", " "),
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
