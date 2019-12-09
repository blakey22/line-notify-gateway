package template

import (
	"testing"
	"strings"
)

func TestTemplate_RenderGithub(t *testing.T) {
	type fields struct {
		basePath string
		locale   string
	}
	type args struct {
		event   string
		payload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "unknown event",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				event:   "non-existed-event",
				payload: []byte(`{"ref":"simple-tag","ref_type":"tag"}`),
			},
			want: "",
			wantErr: true,
		},
		{
			name: "create event (default)",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				event:   "create",
				payload: []byte(`{"ref":"simple-tag","ref_type":"tag"}`),
			},
			want: `
[Event]: Create
[Type]: tag
[Name]: simple-tag`,
			wantErr: false,
		},
		{
			name: "delete event (default)",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				event:   "delete",
				payload: []byte(`{"ref":"simple-tag","ref_type":"tag"}`),
			},
			want: `
[Event]: Delete
[Type]: tag
[Name]: simple-tag`,
			wantErr: false,
		},
		{
			name: "push event (default)",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				event:   "push",
				payload: []byte(`{"ref":"simple-tag","after":"43172584a7e8cfa62b2c9c4aa8036a0a787506b3","compare":"https://github.com/blakey22/hook-test/compare/8f02a3f3ee8d...43172584a7e8"}`),
			},
			want: `
[Event]: Push
[Commit]: 43172584a7e8cfa62b2c9c4aa8036a0a787506b3
[Compare]: https://github.com/blakey22/hook-test/compare/8f02a3f3ee8d...43172584a7e8`,
			wantErr: false,
		},
		
		{
			name: "pull request event (default)",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				event:   "pull_request",
				payload: []byte(`{"action":"open","number":1,"pull_request":{"url":"https://api.github.com/repos/blakey22/hook-test/pulls/1","id":350600950,"node_id":"MDExOlB1bGxSZXF1ZXN0MzUwNjAwOTUw","html_url":"https://github.com/blakey22/hook-test/pull/1","diff_url":"https://github.com/blakey22/hook-test/pull/1.diff","patch_url":"https://github.com/blakey22/hook-test/pull/1.patch","issue_url":"https://api.github.com/repos/blakey22/hook-test/issues/1","number":1,"state":"open"}}`),
			},
			want: `
[Event]: Pull Request
[Status]: open
[URL]: https://github.com/blakey22/hook-test/pull/1`,
			wantErr: false,
		},		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := &Template{
				basePath: tt.fields.basePath,
				locale:   tt.fields.locale,
			}
			got, err := tmpl.RenderGithub(tt.args.event, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.RenderGithub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
				t.Errorf("Template.RenderGithub() = %v, want %v", got, tt.want)
			}
		})
	}
}
