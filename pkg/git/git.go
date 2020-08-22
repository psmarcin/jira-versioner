package git

import (
	"github.com/psmarcin/jira-versioner/pkg/cmd"
	pslog "github.com/psmarcin/jira-versioner/pkg/log"
	"regexp"
)

// Git keeps only dependencies
type Git struct {
	Path         string
	Dependencies Getter
	log          pslog.Logger
}

// Getter is interface for GetTasks dependencies for easier mocking
type Getter interface {
	GetCommits(string, string, string) ([]cmd.Commit, error)
	GetPreviousTag(string, string) (string, error)
}

// New creates Git with default dependencies
func New(path string, log pslog.Logger) Git {
	command := cmd.New(log)
	return Git{
		path,
		command,
		log,
	}
}

// GetTasks gets list of Jira taskIDs from commits
func (g *Git) GetTasks(tag string) ([]string, error) {
	var taskMap = make(map[string]struct{})
	var tasks []string

	previousTag, err := g.Dependencies.GetPreviousTag(tag, g.Path)
	if err != nil {
		return tasks, err
	}
	g.log.Debugf("[GIT] found previous tag: %s", previousTag)

	commits, err := g.Dependencies.GetCommits(tag, previousTag, g.Path)
	if err != nil {
		return nil, err
	}
	g.log.Debugf("[GIT] found commits: %+v", commits)

	re, err := regexp.Compile(`(\w+)-(\d+)`)
	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		issueId := string(re.Find([]byte(commit.Message)))
		if issueId != "" {
			taskMap[issueId] = struct{}{}
		}
	}

	for taskId := range taskMap {
		tasks = append(tasks, taskId)
	}
	g.log.Debugf("[GIT] found tags: %s", tasks)
	return tasks, nil
}
