package dnspod

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupClient() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	params := CommonParams{LoginToken: "dnspod login token"}

	client := NewClient(params)
	client.BaseURL = server.URL + "/"

	return client, mux, func() {
		server.Close()
	}
}

func TestNewClient(t *testing.T) {
	params := CommonParams{LoginToken: "dnspod login token"}

	client := NewClient(params)

	if client.BaseURL != defaultBaseURL {
		t.Errorf("got %v, want %v", client.BaseURL, defaultBaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	params := CommonParams{LoginToken: "dnspod login token"}
	client := NewClient(params)
	client.BaseURL = "https://go.example.com/"

	inPath := "foo"
	req, err := client.NewRequest(http.MethodGet, inPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	// test that relative URL was expanded with the proper BaseURL
	outURL := "https://go.example.com/foo"
	if req.URL.String() != outURL {
		t.Errorf("makeRequest(%v) URL = %v, want %v", inPath, req.URL, outURL)
	}

	// test that default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if client.UserAgent != userAgent {
		t.Errorf("makeRequest() User-Agent = %v, want %v", userAgent, client.UserAgent)
	}
}
