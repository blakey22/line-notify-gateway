package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/pkg/errors"
)

type Template struct {
	basePath string
	locale   string
}

func New() *Template {
	return NewTemplate(flag.Options.Templates, flag.Options.Locale)
}

func NewTemplate(basePath, locale string) *Template {
	return &Template{
		basePath: basePath,
		locale:   locale,
	}
}

func (t *Template) render(name, tmplText string, payload interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(tmplText)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't parse a template for the message")
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, payload)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't execute a template for the message")
	}
	return strings.TrimSpace(buf.String()), nil
}

func (t *Template) load(folder string) (string, error) {
	fn := fmt.Sprintf("%s.tmpl", t.locale)
	path := filepath.Join(t.basePath, folder, fn)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("fail to load template: %s", path)
		return "", err
	}

	return string(dat), nil
}
