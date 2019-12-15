package docker

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blakey22/line-notify-gateway/pkg/handler"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
)

const urlPath string = "/docker-hub"

func init() {
	h := new(dockerHubHandler)
	handler.Register(urlPath, h)
}

type dockerHubHandler struct {
	sendc chan<- line.Notification
	tmpl  *template.Template
}

func (h *dockerHubHandler) Init(sendc chan<- line.Notification) {
	h.sendc = sendc
	h.tmpl = template.New()
}

func (h *dockerHubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail to read body, err: %+v", err)
		handler.WriteJSON(w, http.StatusInternalServerError, "fail to read body")
		return
	}
	defer r.Body.Close()

	msg, err := h.tmpl.RenderDockerHub(b)
	if err != nil {
		log.Printf("fail to process Docker Hub webhook, err: %+v", err)
		handler.WriteJSON(w, http.StatusInternalServerError, "fail to process Docker Hub webhook")
		return
	}

	n := line.Notification{
		Message: msg,
	}
	h.sendc <- n

	handler.WriteJSON(w, http.StatusOK, "ok")
}

func (h *dockerHubHandler) Authorize(r *http.Request) (bool, error) {
	return true, nil
}
