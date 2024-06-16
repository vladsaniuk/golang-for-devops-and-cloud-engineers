package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDoLogin(t *testing.T) {
	var mockToken = Token{
		Token: "456",
	}

	jsonBody, err := json.Marshal(mockToken)
	if err != nil {
		t.Errorf("error marshaling mock token to JSON: %s", err)
	}

	testClient := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: io.NopCloser(bytes.NewBuffer(jsonBody)),
			// Must be set to non-nil value or it panics
			// Header: make(http.Header),
		}
	})

	requestDetails := RequestDetails{
		Token:    "",
		Password: "123",
		URL:      "http://transport-mock/login",
		Client:   *testClient,
	}

	getToken, err := DoLogin(requestDetails)
	if err != nil {
		t.Logf("error making test login request: %s\n", err)
	}

	getToken.GetResponse()

	t.Logf(requestDetails.Token)
	if requestDetails.Token != mockToken.Token {
		t.Error("token is not correct")
	}
}
