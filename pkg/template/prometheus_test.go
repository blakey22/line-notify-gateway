package template

import (
	"strings"
	"testing"
	"time"

	alertmgr "github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
)

func TestTemplate_RenderPrometheus(t *testing.T) {
	type fields struct {
		basePath string
		locale   string
	}
	tests := []struct {
		name    string
		fields  fields
		alert   alertmgr.Alert
		want    string
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			alert: alertmgr.Alert{
				Status: "firing",
				Labels: alertmgr.KV{
					"severity": "line",
				},
				Annotations: alertmgr.KV{
					"summary": "High go threads",
				},
				StartsAt: time.Unix(1575835200, 0),
			},
			want: `
[Status]: firing
[Starts At]: 8 Dec 2019 21:00:00 CET
[Labels]: 
  severity: line
[Annotations]: 
  summary: High go threads
`,
			wantErr: false,
		},
		{
			name: "zh_TW",
			fields: fields{
				basePath: "../../templates",
				locale:   "zh_tw",
			},
			alert: alertmgr.Alert{
				Status: "firing",
				Labels: alertmgr.KV{
					"severity": "line",
				},
				Annotations: alertmgr.KV{
					"summary": "High go threads",
				},
				StartsAt: time.Unix(1575835200, 0),
			},
			want: `
[狀態]: firing
[開始時間]: 8 Dec 2019 21:00:00 CET
[標籤]: 
  severity: line
[註記]: 
  summary: High go threads
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := NewTemplate(tt.fields.basePath, tt.fields.locale)
			got, err := tmpl.RenderPrometheus(tt.alert)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.RenderPrometheus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(got), "message is not matched")
		})
	}
}
