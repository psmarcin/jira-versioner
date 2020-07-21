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

func (m *MockedGit1) GetCommits(tag, prev, gitPath string) ([]cmd.Commit, error) {
	args := m.Called(tag, prev, gitPath)

	return args.Get(0).([]cmd.Commit), args.Error(1)
}

func (m *MockedGit1) GetPreviousTag(tag, gitPath string) (string, error) {
	args := m.Called(tag, gitPath)

	return args.String(0), args.Error(1)
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessage(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae magnam explicabo laborum dolores JIR-15 epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0", ".").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0", ".").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Path:         ".",
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Contains(t, got, "JIR-15")
	assert.Contains(t, got, "JIR-123")
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessageOmitCommitsWithoutTaskID(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae magnam explicabo laborum dolores epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0", ".").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0", ".").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Path:         ".",
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Contains(t, got, "JIR-123")
}

func TestGit_GetTasks_ReturnTaskIDsFromCommitMessageOmitDuplicatedTaskIDs(t *testing.T) {
	firstCommit := cmd.Commit{Hash: "sha1", Message: "feat: JIR-123 Pariatur illum quia nisi praesentium."}
	secondCommit := cmd.Commit{Hash: "sha2", Message: "feat: epudiandae JIR-123 magnam explicabo laborum dolores epudiandae magnam explicabo laborum dolores."}

	m := new(MockedGit1)
	m.On("GetPreviousTag", "v1.1.0", ".").Return("v1.0.0", nil)
	m.On("GetCommits", "v1.1.0", "v1.0.0", ".").Return([]cmd.Commit{
		firstCommit,
		secondCommit,
	}, nil)
	g := &Git{
		Path:         ".",
		Dependencies: m,
	}
	got, err := g.GetTasks("v1.1.0")
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Contains(t, got, "JIR-123")
}
