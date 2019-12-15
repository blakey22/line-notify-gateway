package docker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
)

func Test_dockerHubHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name     string
		event    string
		body     string
		wantCode int
	}{
		{
			name:     "new image",
			body:     `{"callback_url":"https://registry.hub.docker.com/u/svendowideit/testhook/hook/2141b5bi5i5b02bec211i4eeih0242eg11000a/","push_data":{"images":["27d47432a69bca5f2700e4dff7de0388ed65f9d3fb1ec645e2bc24c223dc1cc3","51a9c7c1f8bb2fa19bcd09789a34e63f35abb80044bc10196e304f6634cc582c","..."],"pushed_at":1417566161,"pusher":"trustedbuilder","tag":"latest"},"repository":{"comment_count":0,"date_created":1417494799,"description":"","dockerfile":"#\n# BUILD\t\tdocker build -t svendowideit/apt-cacher .\n# RUN\t\tdocker run -d -p 3142:3142 -name apt-cacher-run apt-cacher\n#\n# and then you can run containers with:\n# \t\tdocker run -t -i -rm -e http_proxy http://192.168.1.2:3142/ debian bash\n#\nFROM\t\tubuntu\n\n\nVOLUME\t\t[/var/cache/apt-cacher-ng]\nRUN\t\tapt-get update ; apt-get install -yq apt-cacher-ng\n\nEXPOSE \t\t3142\nCMD\t\tchmod 777 /var/cache/apt-cacher-ng ; /etc/init.d/apt-cacher-ng start ; tail -f /var/log/apt-cacher-ng/*\n","full_description":"Docker Hub based automated build from a GitHub repo","is_official":false,"is_private":true,"is_trusted":true,"name":"testhook","namespace":"svendowideit","owner":"svendowideit","repo_name":"svendowideit/testhook","repo_url":"https://registry.hub.docker.com/u/svendowideit/testhook/","star_count":0,"status":"Active"}}`,
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid payload",
			body:     `{`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &dockerHubHandler{
				tmpl: template.New(),
			}
			h.Init(make(chan line.Notification, 1))

			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

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

func Test_dockerHubHandler_Authorize(t *testing.T) {
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
			secret:  "secret token",
			body:    "{}",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &dockerHubHandler{}
			req, err := http.NewRequest("POST", urlPath, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			got, err := h.Authorize(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("dockerHubHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dockerHubHandler.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
