package jira

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	pslog "github.com/psmarcin/jira-versioner/pkg/log"
)

// Jira has all necessary details for interacting with Jira service
type Jira struct {
	Client    *jira.Client
	Project   *jira.Project
	ProjectID string
	Version   *jira.Version
	log       pslog.Logger
	dryRun    bool
}

type UpdatePayload struct {
	Update UpdateTypePayload `json:"update"`
}
type UpdateTypePayload struct {
	FixVersions []AddFixedVersion `json:"fixVersions"`
}
type AddFixedVersion struct {
	Add VersionID `json:"add"`
}

type VersionID struct {
	ID string `json:"id"`
}

type Config struct {
	Username  string
	Token     string
	ProjectID string
	BaseURL   string
	Log       pslog.Logger
	DryRun    bool
}

// New creates Jira instance with all required details like email, Token, base url
func New(config *Config) (Jira, error) {
	j := Jira{
		log:    config.Log,
		dryRun: config.DryRun,
	}
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Token,
	}

	client, err := jira.NewClient(tp.Client(), config.BaseURL)
	if err != nil {
		return j, err
	}

	j.Client = client

	_, err = j.getProject(config.ProjectID)
	if err != nil {
		return j, err
	}

	return j, nil
}

// getProject tries to find provided Jira project
func (j *Jira) getProject(projectID string) (jira.Project, error) {
	j.log.Debugf("[JIRA] getting project id from slug: %s", projectID)
	p, _, err := j.Client.Project.Get(projectID)
	if err != nil {
		return jira.Project{}, err
	}
	j.log.Debugf("[JIRA] found project %s", p.Self)

	j.Project = p
	j.ProjectID = p.ID

	j.log.Debugf("[JIRA] project id set to %s", j.ProjectID)

	return *p, nil
}

// GetVersion looks for given version name if exists
func (j Jira) GetVersion(name string) (*jira.Version, bool, error) {
	for i := range j.Project.Versions {
		if j.Project.Versions[i].Name == name {
			return &j.Project.Versions[i], true, nil
		}
	}
	j.log.Debugf("[JIRA] can't find version %s", name)
	return &jira.Version{}, false, nil
}

// CreateVersion creates version in Jira
func (j *Jira) CreateVersion(name string) (*jira.Version, error) {
	version, isFound, err := j.GetVersion(name)
	if err != nil {
		return version, err
	}

	j.log.Debugf("[JIRA] creating version: %s", version.Name)

	if isFound {
		j.Version = version
		j.log.Infof("[JIRA] version %s already exists, skip creating", j.Version.Name)
		return version, nil
	}

	projectID, err := strconv.Atoi(j.ProjectID)
	if err != nil {
		return &jira.Version{}, err
	}

	v := &jira.Version{
		Name:        name,
		ProjectID:   projectID,
		Archived:    false,
		Released:    false,
		StartDate:   time.Now().String(),
		ReleaseDate: time.Now().String(),
		// TODO: put task ids into description
		Description: "",
	}

	if !j.dryRun {
		v, _, err = j.Client.Version.Create(v)
		if err != nil {
			return v, err
		}
	}

	j.Version = v

	j.log.Infof("[JIRA] version created %s", j.Version.Name)

	return version, nil
}

// LinkTasksToVersion iterates over all give tasks and tries to link them to version
func (j Jira) LinkTasksToVersion(taskIds []string) {
	for _, taskID := range taskIds {
		j.log.Debugf("[JIRA] linking %s to %s", taskID, j.Version.Name)

		err := j.SetIssueVersion(taskID)
		if err != nil {
			j.log.Warnf("[JIRA] can't update task %s to fixed version %s (%s)", taskID, j.Version.Name, j.Version.ID)
		}
	}
}

// SetIssueVersion makes http request to Jira service to update task with fixed version
func (j Jira) SetIssueVersion(taskID string) error {
	var res *jira.Response
	p := UpdatePayload{
		Update: UpdateTypePayload{
			FixVersions: []AddFixedVersion{
				{
					Add: VersionID{
						ID: j.Version.ID,
					},
				},
			},
		},
	}

	j.log.Debugf("[JIRA] setting version %s for task %s", j.Version.Name, taskID)
	req, err := j.Client.NewRequest("PUT", "/rest/api/2/issue/"+taskID, p)
	if err != nil {
		return errors.Wrapf(err, "can't create Jira request to %s", "/rest/api/2/issue/"+taskID)
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	if !j.dryRun {
		res, err = j.Client.Do(req, nil)
	}

	if err != nil {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return readErr
		}

		j.log.Warnf("[JIRA] error while setting task %s to %s, %s", taskID, j.Version.Name, body)

		return errors.Wrap(err, string(body))
	}

	j.log.Infof("[JIRA] task updated %s", taskID)
	return nil
}
