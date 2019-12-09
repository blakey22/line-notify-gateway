package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blakey22/line-notify-gateway/pkg/line"
)

type DummyHandler struct {
	sendc chan<- line.Notification
}

func (h *DummyHandler) Init(sendc chan<- line.Notification) {
	h.sendc = sendc
}

func (h *DummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (h DummyHandler) Authorize(r *http.Request) (bool, error) {
	return true, nil
}

func TestRegisterTwice(t *testing.T) {
	Register("/dummy", &DummyHandler{})

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Register twice should panic")
		}
	}()

	Register("/dummy", &DummyHandler{})
}

func TestRegisterNil(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Register nil should panic")
		}
	}()

	Register("/dummy", nil)
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	err := WriteJSON(w, http.StatusOK, "ok")
	if err != nil {
		t.Errorf("fail to write JSON, err: %+v", err)
	}

	expected := `{"status":200,"message":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("WriteJSON writes unexpected body: got %v want %v",
			w.Body.String(),
			expected)
	}

	expected = `application/json`
	if ctype := w.Header().Get("Content-Type"); ctype != expected {
		t.Errorf("WriteJSON writes unexpected Content-Type: got %v want %v",
			ctype,
			expected)
	}
}
