package prometheus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/handler"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/blakey22/line-notify-gateway/pkg/template"
	alertmgr "github.com/prometheus/alertmanager/template"
)

const urlPath string = "/prometheus"

func init() {
	h := new(prometheusHandler)
	handler.Register(urlPath, h)
}

type prometheusHandler struct {
	sendc chan<- line.Notification
	tmpl  *template.Template
}

func (h *prometheusHandler) Init(sendc chan<- line.Notification) {
	h.sendc = sendc
	h.tmpl = template.New()
}

func (h *prometheusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail to read body, err: %+v", err)
		handler.WriteJSON(w, http.StatusInternalServerError, "fail to read body")
		return
	}
	defer r.Body.Close()

	data := alertmgr.Data{}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Printf("unable to decode body, err: %+v", err)
		handler.WriteJSON(w, http.StatusBadRequest, "unable to decode body")
		return
	}

	for _, alert := range data.Alerts {
		msg, err := h.tmpl.RenderPrometheus(alert)
		if err != nil {
			log.Printf("fail to parse prometheus alert, err: %+v", err)
			handler.WriteJSON(w, http.StatusInternalServerError, "fail to parse prometheus alert")
			return
		}
		n := line.Notification{
			Message: msg,
		}
		h.sendc <- n
	}

	handler.WriteJSON(w, http.StatusOK, "ok")
}

func (h prometheusHandler) Authorize(r *http.Request) (bool, error) {
	if len(flag.Options.Secret) == 0 {
		return true, nil
	}

	if r.Header.Get("Authorization") == "Bearer "+flag.Options.Secret {
		return true, nil
	}

	return false, nil
}
