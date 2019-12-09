package github

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
)

func Test_githubHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name     string
		event    string
		body     string
		wantCode int
	}{
		{
			name:     "ping event",
			event:    "ping",
			body:     "",
			wantCode: http.StatusOK,
		},
		{
			name:     "create event with secret",
			event:    "create",
			body:     "{}",
			wantCode: http.StatusOK,
		},
		{
			name:     "unknown event",
			event:    "not-existed",
			body:     "{}",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &githubHandler{
				tmpl: template.New(),
			}
			h.Init(make(chan line.Notification, 1))

			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("X-GitHub-Event", tt.event)

			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("handler returned unexpected status code: got %v want %v",
					w.Code,
					tt.wantCode)
			}
		})
	}
}

func Test_githubHandler_Authorize(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		sign    string
		body    string
		want    bool
		wantErr bool
	}{
		{
			name:    "no secret",
			secret:  "",
			sign:    "",
			body:    "{}",
			want:    true,
			wantErr: false,
		},
		{
			name:    "has secret",
			secret:  "d50a98b0-34f7-4de2-a0d4-3d939705396e",
			sign:    "sha1=3a00cf0b42f7dac4d3a8d368bb1a0be221ba0559",
			body:    "{}",
			want:    true,
			wantErr: false,
		},
		{
			name:    "has secret and wrong sign",
			secret:  "d50a98b0-34f7-4de2-a0d4-3d939705396e",
			sign:    "wrong sign",
			body:    "{}",
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.Options.Secret = tt.secret
			h := &githubHandler{}
			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("X-Hub-Signature", tt.sign)

			got, err := h.Authorize(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("githubHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("githubHandler.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
