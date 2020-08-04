package jira

import (
	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

// Jira has all necessary details for interacting with Jira service
type Jira struct {
	token     string
	Client    *jira.Client
	Project   *jira.Project
	ProjectID string
	Version   *jira.Version
}

type UpdatePayload struct {
	Update UpdateTypePayload `json:"update"`
}
type UpdateTypePayload struct {
	FixVersions []AddFixedVersion `json:"fixVersions"`
}
type AddFixedVersion struct {
	Add IdVersion `json:"add"`
}

type IdVersion struct {
	Id string `json:"id"`
}

// New creates Jira instance with all required details like email, token, base url
func New(email, token, projectId, baseUrl string) (Jira, error) {
	j := Jira{}
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return j, err
	}

	j.Client = client

	_, err = j.getProject(projectId)
	if err != nil {
		return j, err
	}

	return j, nil
}

// getProject tries to find provided Jira project
func (j *Jira) getProject(projectId string) (jira.Project, error) {
	p, _, err := j.Client.Project.Get(projectId)
	if err != nil {
		return jira.Project{}, err
	}
	j.Project = p
	j.ProjectID = projectId

	return *p, nil
}

// GetVersion looks for given version name if exists
func (j Jira) GetVersion(name string) (*jira.Version, bool, error) {
	for _, version := range j.Project.Versions {
		if version.Name == name {
			return &version, true, nil
		}
	}
	return &jira.Version{}, false, nil
}

// CreateVersion creates version in Jira
func (j *Jira) CreateVersion(name string) (*jira.Version, error) {
	version, isFound, err := j.GetVersion(name)
	if err != nil {
		return version, err
	}
	if isFound == true {
		j.Version = version
		log.Printf("[JIRA] version %s already exists, skip creating", j.Version.Name)
		return version, nil
	}

	projectId, err := strconv.Atoi(j.ProjectID)
	if err != nil {
		return &jira.Version{}, err
	}

	v := &jira.Version{
		Name:        name,
		ProjectID:   projectId,
		Archived:    false,
		Released:    false,
		StartDate:   time.Now().String(),
		ReleaseDate: time.Now().String(),
		// TODO: put task ids into description
		Description: "",
	}
	version, _, err = j.Client.Version.Create(v)
	if err != nil {
		return v, err
	}

	j.Version = version

	log.Printf("[JIRA] version created %s", j.Version.Name)

	return version, nil
}

// LinkTasksToVersion iterates over all give tasks and tries to link them to version
func (j Jira) LinkTasksToVersion(taskIds []string) {
	for _, taskId := range taskIds {
		err := j.SetIssueVersion(taskId)
		if err != nil {
			log.Printf("[JIRA] can't update task %s to fixed version %s (%s)", taskId, j.Version.Name, j.Version.ID)
		}
	}
}

// SetIssueVersion makes http request to Jira service to update task with fixed version
func (j Jira) SetIssueVersion(taskID string) error {
	p := UpdatePayload{
		Update: UpdateTypePayload{
			FixVersions: []AddFixedVersion{
				{
					Add: IdVersion{
						Id: j.Version.ID,
					},
				},
			},
		},
	}

	req, _ := j.Client.NewRequest("PUT", "/rest/api/2/issue/"+taskID, p)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := j.Client.Do(req, nil)
	if err != nil {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return readErr
		}

		return errors.Wrap(err, string(body))
	}

	log.Printf("[JIRA] task updated %s", taskID)
	return nil
}
