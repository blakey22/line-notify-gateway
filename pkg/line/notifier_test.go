package line

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifier_send(t *testing.T) {
	token := "3939889"
	message := "test failed"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+token)
		assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		assert.NoError(t, req.ParseForm(), "fail to parse post data")
		assert.Equal(t, req.Form.Get("message"), message)
		
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	n := Notification{
		Message: message,
	}
	notifier := NewNotifier(server.URL, token)
	assert.NoError(t, notifier.send(n), "send failed")
}
