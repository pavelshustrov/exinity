package gateway_a

import (
	"context"
	"exinity/internal/clients"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientWithResponses_Call_Success(t *testing.T) {
	// Create a test server that responds with a successful response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/payments", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`)) // Mocked response body
	}))
	defer ts.Close()

	// Mock the ClientWithResponses and call the function
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // Set the timeout for the request
		Transport: &clients.CustomTransport{Transport: http.DefaultTransport},
	}

	client, err := NewClientWithResponses(ts.URL, WithHTTPClient(httpClient))
	assert.NoError(t, err)

	err = client.Call(context.Background(), "txn123", 100.50, "USD")

	assert.NoError(t, err, "Expected no error for successful API call")
}

func TestClientWithResponses_Call_Failure(t *testing.T) {
	// Create a test server that responds with an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Mock the ClientWithResponses and call the function
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // Set the timeout for the request
		Transport: &clients.CustomTransport{Transport: http.DefaultTransport},
	}

	client, err := NewClientWithResponses(ts.URL, WithHTTPClient(httpClient))
	assert.NoError(t, err)

	err = client.Call(context.Background(), "txn123", 100.50, "USD")

	assert.Error(t, err, "Expected an error for failed API call")
}

func TestClientWithResponses_Call_Timeout(t *testing.T) {
	// Create a test server that simulates a timeout
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second)
	}))
	defer ts.Close()

	// Mock the ClientWithResponses and call the function
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // Set the timeout for the request
		Transport: &clients.CustomTransport{Transport: http.DefaultTransport},
	}

	client, err := NewClientWithResponses(ts.URL, WithHTTPClient(httpClient))
	assert.NoError(t, err)

	err = client.Call(context.Background(), "txn123", 100.50, "USD")

	assert.Error(t, err, "Expected a timeout error")
	assert.Contains(t, err.Error(), "Client.Timeout exceeded", "Expected timeout error message")
}
