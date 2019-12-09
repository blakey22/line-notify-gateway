package template

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// RenderGithub render Github event to localized content
func (t *Template) RenderGithub(event string, payload []byte) (string, error) {
	switch event {
	case "create":
		return t.githubCreate(payload)
	case "delete":
		return t.githubDelete(payload)
	case "push":
		return t.githubPush(payload)
	case "pull_request":
		return t.githubPullRequest(payload)
	default:
		return "", errors.Errorf("unkonwn github event: %s", event)
	}
}

func (t *Template) githubCreate(payload []byte) (string, error) {
	const defaultTemplate = `
[Event]: Create
[Type]: {{ .RefType }}
[Name]: {{ .Ref }}
	`

	var data struct {
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	tmpl, err := t.load("github/create")
	if err != nil {
		tmpl = defaultTemplate
	}

	return t.render("github-create", tmpl, data)
}

func (t *Template) githubDelete(payload []byte) (string, error) {
	const defaultTemplate = `
[Event]: Delete
[Type]: {{ .RefType }}
[Name]: {{ .Ref }}
	`

	var data struct {
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	tmpl, err := t.load("github/delete")
	if err != nil {
		tmpl = defaultTemplate
	}

	return t.render("github-delete", tmpl, data)
}

func (t *Template) githubPush(payload []byte) (string, error) {
	const defaultTemplate = `
[Event]: Push
[Commit]: {{ .After }}
[Compare]: {{ .Compare }}
	`

	var data struct {
		Ref     string `json:"ref"`
		After   string `json:"after"`
		Compare string `json:"compare"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	tmpl, err := t.load("github/push")
	if err != nil {
		tmpl = defaultTemplate
	}

	return t.render("github-push", tmpl, data)
}

func (t *Template) githubPullRequest(payload []byte) (string, error) {
	const defaultTemplate = `
[Event]: Pull Request
[Status]: {{ .PullRequest.State }}
[URL]: {{ .PullRequest.HTMLURL }}
	`

	var data struct {
		Action      string `json:"action"`
		Number      int    `json:"number"`
		PullRequest struct {
			URL      string `json:"url"`
			ID       int    `json:"id"`
			NodeID   string `json:"node_id"`
			HTMLURL  string `json:"html_url"`
			DiffURL  string `json:"diff_url"`
			PatchURL string `json:"patch_url"`
			IssueURL string `json:"issue_url"`
			Number   int    `json:"number"`
			State    string `json:"state"`
			Locked   bool   `json:"locked"`
			Title    string `json:"title"`
		} `json:"pull_request"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return "", err
	}

	tmpl, err := t.load("github/pull-request")
	if err != nil {
		tmpl = defaultTemplate
	}

	return t.render("github-pull-request", tmpl, data)
}
