package prometheus

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
)

var message string = `
{
	"receiver": "line-notify",
	"status": "firing",
	"alerts": [
	  {
		"status": "firing",
		"labels": {
		  "alertname": "HighGoThreads",
		  "instance": "localhost:9090",
		  "job": "prometheus",
		  "severity": "line"
		},
		"annotations": {
		  "summary": "High go threads"
		},
		"startsAt": "2019-12-08T14:00:46.431338009Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://f12169ece5f7:9090/graph?g0.expr=go_threads+%3E+5&g0.tab=1",
		"fingerprint": "471410d1530c7fa1"
	  }
	],
	"groupLabels": {
	  "alertname": "HighGoThreads",
	  "instance": "localhost:9090",
	  "job": "prometheus",
	  "severity": "line"
	},
	"commonLabels": {
	  "alertname": "HighGoThreads",
	  "instance": "localhost:9090",
	  "job": "prometheus",
	  "severity": "line"
	},
	"commonAnnotations": {
	  "summary": "High go threads"
	},
	"externalURL": "http://04018f184e8b:9093",
	"version": "4",
	"groupKey": "{}:{alertname=\"HighGoThreads\", instance=\"localhost:9090\", job=\"prometheus\", severity=\"line\"}"
  }
`

func Test_prometheusHandler_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", urlPath, strings.NewReader(message))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	h := &prometheusHandler{
		tmpl: template.New(),
	}
	h.Init(make(chan line.Notification, 1))

	h.ServeHTTP(w, req)

	expected := `{"status":200,"message":"ok"}`
	if body := w.Body.String(); body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body,
			expected)
	}

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned unexpected status: got %v want %v",
			status,
			expected)
	}
}

func Test_prometheusHandler_Authorize(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		body    string
		want    bool
		wantErr bool
	}{
		{
			name:    "no secret",
			secret:  "",
			body:    "{}",
			want:    true,
			wantErr: false,
		},
		{
			name:    "has secret",
			secret:  "d50a98b0-34f7-4de2-a0d4-3d939705396e",
			body:    "{}",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.Options.Secret = tt.secret
			h := &prometheusHandler{}
			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Authorization", "Bearer "+tt.secret)

			got, err := h.Authorize(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("prometheusHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("prometheusHandler.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
