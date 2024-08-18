package clients

import (
	"github.com/avast/retry-go"
	"net/http"
	"time"
)

// CustomTransport is an HTTP transport that adds retries.
type CustomTransport struct {
	Transport http.RoundTripper
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	err := retry.Do(
		func() error {
			var err error
			resp, err = t.Transport.RoundTrip(req)
			if err != nil || resp.StatusCode >= 500 {
				// Consider retrying only on server errors (5xx)
				return err
			}
			return nil
		},
		retry.Attempts(3),          // Number of retries
		retry.Delay(2*time.Second), // Delay between retries
	)
	return resp, err
}
