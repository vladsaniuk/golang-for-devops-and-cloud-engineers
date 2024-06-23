package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// create mock HTTP response via Client and Round Tripper
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
			Body:       io.NopCloser(bytes.NewBuffer(jsonBody)),
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
		t.Errorf("error making test login request: %s\n", err)
	}

	requestDetails.Token = getToken.GetResponse()

	t.Logf("requestDetails.Token is %v\n", requestDetails.Token)
	if requestDetails.Token != mockToken.Token {
		t.Error("token is not correct")
	}
}
