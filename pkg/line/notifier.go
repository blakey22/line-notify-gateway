package line

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Notifier struct {
	endpoint string
	token    string
	client   *http.Client
}

func NewNotifier(endpoint, token string) *Notifier {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return &Notifier{
		endpoint: endpoint,
		token:    token,
		client:   client,
	}
}

func (nt *Notifier) Run(ctx context.Context, sendc <-chan Notification) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case n := <-sendc:
			if err := nt.send(n); err != nil {
				return err
			}
		}
	}
}

func (nt *Notifier) send(n Notification) error {
	body := n.URLValues()
	req, err := http.NewRequest("POST", nt.endpoint, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+nt.token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := nt.client.Do(req)
	if err != nil {
		return errors.Errorf("unable to send the line notification, err: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Errorf("line server responded an error: %s", content)
		}
	}

	return nil
}
