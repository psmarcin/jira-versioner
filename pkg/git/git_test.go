package git

import (
	"github.com/psmarcin/jira-versioner/pkg/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockedGit1 struct {
	mock.Mock
}

func (m *MockedGit1) GetCommits(tag, prev string) ([]cmd.Commit, error){
	args := m.Called(tag, prev)

	return args.Get(0).([]cmd.Commit), args.Error(1)
}

func (m *MockedGit1) GetPreviousTag(tag string) (string, error){
	args := m.Called(tag)

	return args.String(0), args.Error(1)
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessage(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae magnam explicabo laborum dolores JIR-15 epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got,2)
	assert.EqualValues(t, got[0], "JIR-123")
	assert.EqualValues(t, got[1], "JIR-15")
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessageOmitCommitsWithoutTaskID(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae magnam explicabo laborum dolores epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got,1)
	assert.EqualValues(t, got[0], "JIR-123")
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessageOmitDuplicatedTaskIDs(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae JIR-123 magnam explicabo laborum dolores epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got,1)
	assert.EqualValues(t, got[0], "JIR-123")
}
