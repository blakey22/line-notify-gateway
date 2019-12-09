package template

import (
	"github.com/prometheus/alertmanager/template"
)

// RenderPrometheus render Alert from Prometheus Alert manager to localized content
func (t *Template) RenderPrometheus(alert template.Alert) (string, error) {
	const defaultTmpl = `
[Status]: {{ .Status }}
[Starts At]: {{ .StartsAt.Format "2 Jan 2006 15:04:05 MST" }}
[Labels]: {{ range $key, $value := .Labels }}
  {{ $key }}: {{ $value }}{{ end }}
[Annotations]: {{ range $key, $value := .Annotations }}
  {{ $key }}: {{ $value }}{{ end }}
`

	tmpl, err := t.load("prometheus")
	if err != nil {
		tmpl = defaultTmpl
	}

	content, err := t.render("prometheus", tmpl, alert)

	return content, err
}
