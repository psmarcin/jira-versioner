package cmd

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestGitCommand_GetPreviousTag(t *testing.T) {
	type fields struct {
		PreviousTagGetter PreviousTagGetter
		CommitGetter      CommitGetter
	}
	type args struct {
		tag           string
		prevTagGetter PreviousTagGetter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return v1.0.0",
			fields: fields{
				PreviousTagGetter: func(name string, arg ...string) (string, error) {
					return "v1.0.0", nil
				},
				CommitGetter: nil,
			},
			args: args{
				tag: "v1.1.0",
			},
			want:    "v1.0.0",
			wantErr: false,
		},
		{
			name: "should return error",
			fields: fields{
				PreviousTagGetter: func(name string, arg ...string) (string, error) {
					return "v1.0.0", errors.New("err 128")
				},
				CommitGetter: nil,
			},
			args: args{
				tag: "v1.1.0",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should trim whitespace in result",
			fields: fields{
				PreviousTagGetter: func(name string, arg ...string) (string, error) {
					return `     v1.0.0
`, nil
				},
				CommitGetter: nil,
			},
			args: args{
				tag: "v1.1.0",
			},
			want:    "v1.0.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Git{
				PreviousTagGetter: tt.fields.PreviousTagGetter,
				CommitGetter:      tt.fields.CommitGetter,
			}
			got, err := c.GetPreviousTag(tt.args.tag, ".")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPreviousTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPreviousTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitCommand_GetCommits(t *testing.T) {
	type fields struct {
		PreviousTagGetter PreviousTagGetter
		CommitGetter      CommitGetter
	}
	type args struct {
		tag         string
		previousTag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Commit
		wantErr bool
	}{
		{
			name: "should return two commits from v1.1.0 to v1.0.0",
			fields: fields{
				PreviousTagGetter: nil,
				CommitGetter: func(name string, arg ...string) (string, error) {
					return `sha1;feat: JIR-1556 commit message
sha2;fix: JIR-9899 commit message`, nil
				},
			},
			args: args{
				tag:         "v1.1.0",
				previousTag: "v1.0.0",
			},
			want: []Commit{
				{
					Hash:    "sha1",
					Message: "feat: JIR-1556 commit message",
				},
				{
					Hash:    "sha2",
					Message: "fix: JIR-9899 commit message",
				},
			},
			wantErr: false,
		},
		{
			name: "should return no commits",
			fields: fields{
				PreviousTagGetter: nil,
				CommitGetter: func(name string, arg ...string) (string, error) {
					return ``, nil
				},
			},
			args: args{
				tag:         "v1.1.0",
				previousTag: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "should return error from command",
			fields: fields{
				PreviousTagGetter: nil,
				CommitGetter: func(name string, arg ...string) (string, error) {
					return ``, errors.New("err 128")
				},
			},
			args: args{
				tag:         "v1.1.0",
				previousTag: "v1.0.0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := zap.NewExample().Sugar()
			defer log.Sync()
			c := Git{
				PreviousTagGetter: tt.fields.PreviousTagGetter,
				CommitGetter:      tt.fields.CommitGetter,
				log:               log,
			}
			got, err := c.GetCommits(tt.args.tag, tt.args.previousTag, ".")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	log := zap.NewExample().Sugar()
	defer log.Sync()
	g := New(log)
	assert.NotEmpty(t, g)
}
