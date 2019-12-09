package tester

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/line"
)

func Test_testerHandler_ServeHTTP(t *testing.T) {
	message := `test`
	req, err := http.NewRequest("POST", urlPath, strings.NewReader(message))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	sendc := make(chan line.Notification, 1)
	h := new(testerHandler)
	h.Init(sendc)
	h.ServeHTTP(w, req)
	close(sendc)

	expected := `{"status":200,"message":"test"}`
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

	for n := range sendc {
		if n.Message != message {
			t.Errorf("send channel got unexpected status: got %v want %v",
				n.Message,
				message)
		}
	}
}

func Test_testerHandler_Authorize(t *testing.T) {
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
			h := &testerHandler{}
			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Authorization", "Bearer "+tt.secret)

			got, err := h.Authorize(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("testerHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("testerHandler.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
