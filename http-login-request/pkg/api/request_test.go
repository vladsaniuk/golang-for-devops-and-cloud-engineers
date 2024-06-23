package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestDoRequestWords(t *testing.T) {
	words := Words{
		Page:  "words",
		Input: "",
		Words: []string{
			"word1",
			"word2",
		},
	}

	jsonBody, err := json.Marshal(words)
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
		Token:    "456",
		Password: "123",
		URL:      "http://transport-mock/words",
		Client:   *testClient,
	}

	response, err := DoRequest(requestDetails)
	slog.Debug("main response: " + response.GetResponse())
	if err != nil {
		slog.Error("error making request: " + err.Error())
		os.Exit(1)
	} else if response == nil {
		slog.Error("Something went wrong - got nil in response")
		os.Exit(1)
	}

	responseResult := response.GetResponse()

	t.Logf("responseResult is %v\n", responseResult)

	if responseResult != fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\nInput: %s\n", words.Page, strings.Join(words.Words, ", "), words.Input) {
		t.Error("responseResult is not the same as response mock")
	}
}
