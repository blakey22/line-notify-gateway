package tester

import (
	"io/ioutil"
	"net/http"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/handler"
	"github.com/blakey22/line-notify-gateway/pkg/line"
)

const urlPath string = "/tester"

func init() {
	h := new(testerHandler)
	handler.Register(urlPath, h)
}

type testerHandler struct {
	sendc chan<- line.Notification
}

func (h *testerHandler) Init(sendc chan<- line.Notification) {
	h.sendc = sendc
}

func (h *testerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handler.WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	message := string(data)
	n := line.Notification{
		Message: message,
	}
	h.sendc <- n
	handler.WriteJSON(w, http.StatusOK, message)
}

func (h testerHandler) Authorize(r *http.Request) (bool, error) {
	if len(flag.Options.Secret) == 0 {
		return true, nil
	}

	if r.Header.Get("Authorization") == "Bearer "+flag.Options.Secret {
		return true, nil
	}

	return false, nil
}
