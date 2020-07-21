package cmd

import "testing"

func TestExec(t *testing.T) {
	type args struct {
		name string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return echo",
			args: args{
				name: "echo",
				args: []string{"123"},
			},
			want: `123
`,
			wantErr: false,
		},
		{
			name: "should return error",
			args: args{
				name: "asdasda", // command doesn't exists
				args: nil,
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
