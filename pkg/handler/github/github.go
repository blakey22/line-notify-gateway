package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/handler"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
)

const urlPath string = "/github"

func init() {
	h := new(githubHandler)
	handler.Register(urlPath, h)
}

type githubHandler struct {
	sendc chan<- line.Notification
	tmpl  *template.Template
}

func (h *githubHandler) Init(sendc chan<- line.Notification) {
	h.sendc = sendc
	h.tmpl = template.New()
}

func (h *githubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-GitHub-Event")
	if event == "ping" {
		handler.WriteJSON(w, http.StatusOK, "ok")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail to read body, err: %+v", err)
		handler.WriteJSON(w, http.StatusInternalServerError, "fail to read body")
		return
	}
	defer r.Body.Close()

	msg, err := h.tmpl.RenderGithub(event, b)
	if err != nil {
		log.Printf("fail to process Github event, err: %+v", err)
		handler.WriteJSON(w, http.StatusInternalServerError, "fail to process Github event")
		return
	}

	n := line.Notification{
		Message: msg,
	}
	h.sendc <- n

	handler.WriteJSON(w, http.StatusOK, "ok")
}

func (h *githubHandler) Authorize(r *http.Request) (bool, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	r.Body.Close()

	// restore Body state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	secret := flag.Options.Secret
	sign := r.Header.Get("X-Hub-Signature")
	if len(secret) == 0 {
		return true, nil
	}

	sha := hmac.New(sha1.New, []byte(secret))
	sha.Write(payload)
	calcSign := fmt.Sprintf("sha1=%s", hex.EncodeToString(sha.Sum(nil)))
	return calcSign == sign, nil
}
