package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blakey22/line-notify-gateway/pkg/line"
)

var handlersMu sync.Mutex
var handlers = make(map[string]Handler)

type Handler interface {
	Init(sendc chan<- line.Notification)
	Authorize(r *http.Request) (bool, error)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func Register(path string, handler Handler) {
	handlersMu.Lock()
	defer handlersMu.Unlock()

	if handler == nil {
		panic("Register handler is nil")
	}

	if _, existed := handlers[path]; existed {
		panic("Register the same path twice: " + path)
	}
	
	handlers[path] = handler
}

// Count provides the count of registed handlers
func Count() int {
	handlersMu.Lock()
	defer handlersMu.Unlock()

	return len(handlers)
}

// Setup initialize handlers and register http paths
func Setup(sendc chan<- line.Notification) {
	for path, handler := range handlers {
		handler.Init(sendc)
		http.HandleFunc(path, setupHandler(handler))
	}
}

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteJSON writes status and message into JSON response
func WriteJSON(w http.ResponseWriter, status int, message string) error {
	data := response{
		Status:  status,
		Message: message,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, string(bytes))

	return nil
}

func setupHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, err := handler.Authorize(r)
		if err != nil {
			log.Printf("unable to authorize, err: %+v", err)
			WriteJSON(w, http.StatusInternalServerError, "unable to authorize")
			return
		}
		if !auth {
			log.Printf("authorize failed, err: %+v", err)
			WriteJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		handler.ServeHTTP(w, r)
	}
}
