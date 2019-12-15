package template

import (
	"encoding/json"
)

// RenderDockerHub render Github event to localized content
func (t *Template) RenderDockerHub(payload []byte) (string, error) {
	const defaultTemplate = `
[Event]: New Docker Hub Image
[Name]: {{ .Repository.RepoName }}:{{ .PushData.Tag }}
	`

	var data struct {
		CallbackURL string `json:"callback_url"`
		PushData    struct {
			Images   []string `json:"images"`
			PushedAt float64  `json:"pushed_at"`
			Pusher   string   `json:"pusher"`
			Tag      string   `json:"tag"`
		} `json:"push_data"`
		Repository struct {
			CommentCount    int     `json:"comment_count"`
			DateCreated     float64 `json:"date_created"`
			Description     string  `json:"description"`
			Dockerfile      string  `json:"dockerfile"`
			FullDescription string  `json:"full_description"`
			IsOfficial      bool    `json:"is_official"`
			IsPrivate       bool    `json:"is_private"`
			IsTrusted       bool    `json:"is_trusted"`
			Name            string  `json:"name"`
			Namespace       string  `json:"namespace"`
			Owner           string  `json:"owner"`
			RepoName        string  `json:"repo_name"`
			RepoURL         string  `json:"repo_url"`
			StarCount       int     `json:"star_count"`
			Status          string  `json:"status"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	tmpl, err := t.load("docker-hub")
	if err != nil {
		tmpl = defaultTemplate
	}

	return t.render("docker-hub", tmpl, data)
}
