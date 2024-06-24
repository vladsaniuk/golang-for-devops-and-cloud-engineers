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
		t.Errorf("error marshaling mock Words to JSON: %s", err)
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
	t.Logf("responseResult is \n%v\n", responseResult)

	responseExpected := fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\nInput: %s\n", words.Page, strings.Join(words.Words, ", "), words.Input)
	t.Logf("responseExpected is \n%v\n", responseExpected)

	if responseResult != responseExpected {
		t.Error("responseResult is not the same as response mock")
	}
}

func TestDoRequestOccurrence(t *testing.T) {
	occurrence := Occurrence{
		Page: "occurrence",
		Words: map[string]int{
			"word2": 32,
		},
	}

	jsonBody, err := json.Marshal(occurrence)
	if err != nil {
		t.Errorf("error marshaling mock Occurrence to JSON: %s", err)
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
	t.Logf("responseResult is \n%v\n", responseResult)

	// mimic same output as GetResponse() to compare
	words := []string{}
	for word, occurrence := range occurrence.Words {
		words = append(words, fmt.Sprintf("Word is %s, it showed up %v time", word, occurrence))
	}

	responseExpected := fmt.Sprintf("Parsed JSON:\nPage: %s\nWords:\n%s\n", occurrence.Page, strings.Join(words, "\n"))
	t.Logf("responseExpected is \n%v\n", responseExpected)

	if responseResult != responseExpected {
		t.Error("responseResult is not the same as response mock")
	}
}

func TestDoRequestAssignment1(t *testing.T) {
	assignment1 := Assignment1{
		Page: "assignment1",
		Words: []string{
			"six",
			"two",
			"one",
		},
		Percentages: map[string]float32{
			"one": 0.33,
		},
		Special: []interface{}{
			"one",
			"two",
			nil,
		},
		ExtraSpecial: []interface{}{
			1.0,
			2.0,
			"3",
		},
	}

	jsonBody, err := json.Marshal(assignment1)
	if err != nil {
		t.Errorf("error marshaling mock Occurrence to JSON: %s", err)
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
	t.Logf("responseResult is \n%v\n", responseResult)

	// mimic same output as GetResponse() to compare
	percentages := []string{}
	for number, percentage := range assignment1.Percentages {
		percentages = append(percentages, fmt.Sprintf("Number is %s, it's percentage is %v\n", number, percentage))
	}

	specials := []string{}
	for _, special := range assignment1.Special {
		specials = append(specials, fmt.Sprintf("%v with type of %T\n", special, special))
	}

	extraSpecials := []string{}
	for _, extraSpecial := range assignment1.ExtraSpecial {
		extraSpecials = append(extraSpecials, fmt.Sprintf("%v with type of %T\n", extraSpecial, extraSpecial))
	}

	responseExpected := fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\nPercentages:\n%vSpecials:\n%vExtraSpecials:\n%v", assignment1.Page, strings.Join(assignment1.Words, ", "), strings.Join(percentages, ""), strings.Join(specials, ""), strings.Join(extraSpecials, ""))
	t.Logf("responseExpected is \n%v\n", responseExpected)

	if responseResult != responseExpected {
		t.Error("responseResult is not the same as response mock")
	}
}
