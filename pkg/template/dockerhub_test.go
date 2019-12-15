package template

import (
	"strings"
	"testing"
)

func TestTemplate_RenderDockerHub(t *testing.T) {
	type fields struct {
		basePath string
		locale   string
	}
	type args struct {
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
			name: "default",
			fields: fields{
				basePath: "",
				locale:   "",
			},
			args: args{
				payload: []byte(`{"callback_url":"https://registry.hub.docker.com/u/svendowideit/testhook/hook/2141b5bi5i5b02bec211i4eeih0242eg11000a/","push_data":{"images":["27d47432a69bca5f2700e4dff7de0388ed65f9d3fb1ec645e2bc24c223dc1cc3","51a9c7c1f8bb2fa19bcd09789a34e63f35abb80044bc10196e304f6634cc582c","..."],"pushed_at":1417566161,"pusher":"trustedbuilder","tag":"latest"},"repository":{"comment_count":0,"date_created":1417494799,"description":"","dockerfile":"#\n# BUILD\t\tdocker build -t svendowideit/apt-cacher .\n# RUN\t\tdocker run -d -p 3142:3142 -name apt-cacher-run apt-cacher\n#\n# and then you can run containers with:\n# \t\tdocker run -t -i -rm -e http_proxy http://192.168.1.2:3142/ debian bash\n#\nFROM\t\tubuntu\n\n\nVOLUME\t\t[/var/cache/apt-cacher-ng]\nRUN\t\tapt-get update ; apt-get install -yq apt-cacher-ng\n\nEXPOSE \t\t3142\nCMD\t\tchmod 777 /var/cache/apt-cacher-ng ; /etc/init.d/apt-cacher-ng start ; tail -f /var/log/apt-cacher-ng/*\n","full_description":"Docker Hub based automated build from a GitHub repo","is_official":false,"is_private":true,"is_trusted":true,"name":"testhook","namespace":"svendowideit","owner":"svendowideit","repo_name":"svendowideit/testhook","repo_url":"https://registry.hub.docker.com/u/svendowideit/testhook/","star_count":0,"status":"Active"}}`),
			},
			want: `
[Event]: New Docker Hub Image
[Name]: svendowideit/testhook:latest`,
			wantErr: false,
		},
		{
			name: "zh-TW",
			fields: fields{
				basePath: "../../templates",
				locale:   "zh_TW",
			},
			args: args{
				payload: []byte(`{"callback_url":"https://registry.hub.docker.com/u/svendowideit/testhook/hook/2141b5bi5i5b02bec211i4eeih0242eg11000a/","push_data":{"images":["27d47432a69bca5f2700e4dff7de0388ed65f9d3fb1ec645e2bc24c223dc1cc3","51a9c7c1f8bb2fa19bcd09789a34e63f35abb80044bc10196e304f6634cc582c","..."],"pushed_at":1417566161,"pusher":"trustedbuilder","tag":"latest"},"repository":{"comment_count":0,"date_created":1417494799,"description":"","dockerfile":"#\n# BUILD\t\tdocker build -t svendowideit/apt-cacher .\n# RUN\t\tdocker run -d -p 3142:3142 -name apt-cacher-run apt-cacher\n#\n# and then you can run containers with:\n# \t\tdocker run -t -i -rm -e http_proxy http://192.168.1.2:3142/ debian bash\n#\nFROM\t\tubuntu\n\n\nVOLUME\t\t[/var/cache/apt-cacher-ng]\nRUN\t\tapt-get update ; apt-get install -yq apt-cacher-ng\n\nEXPOSE \t\t3142\nCMD\t\tchmod 777 /var/cache/apt-cacher-ng ; /etc/init.d/apt-cacher-ng start ; tail -f /var/log/apt-cacher-ng/*\n","full_description":"Docker Hub based automated build from a GitHub repo","is_official":false,"is_private":true,"is_trusted":true,"name":"testhook","namespace":"svendowideit","owner":"svendowideit","repo_name":"svendowideit/testhook","repo_url":"https://registry.hub.docker.com/u/svendowideit/testhook/","star_count":0,"status":"Active"}}`),
			},
			want: `
[事件]: 新 Docker Hub Image
[名稱]: svendowideit/testhook:latest`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := &Template{
				basePath: tt.fields.basePath,
				locale:   tt.fields.locale,
			}
			got, err := tmpl.RenderDockerHub(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.RenderDockerHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
				t.Errorf("Template.RenderDockerHub() = %v, want %v", got, tt.want)
			}
		})
	}
}
