package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Assignment1 struct {
	Page         string             `json:"page"`
	Words        []string           `json:"words"`
	Percentages  map[string]float32 `json:"percentages"`
	Special      []interface{}      `json:"special"`
	ExtraSpecial []interface{}      `json:"extraSpecial"`
}

func (p Assignment1) GetResponse() string {
	percentages := []string{}
	for number, percentage := range p.Percentages {
		percentages = append(percentages, fmt.Sprintf("Number is %s, it's percentage is %v\n", number, percentage))
	}

	specials := []string{}
	for _, special := range p.Special {
		specials = append(specials, fmt.Sprintf("%v with type of %T\n", special, special))
	}

	extraSpecials := []string{}
	for _, extraSpecial := range p.ExtraSpecial {
		extraSpecials = append(extraSpecials, fmt.Sprintf("%v with type of %T\n", extraSpecial, extraSpecial))
	}

	return fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\nPercentages:\n%vSpecials:\n%vExtraSpecials:\n%v", p.Page, strings.Join(p.Words, ", "), strings.Join(percentages, ""), strings.Join(specials, ""), strings.Join(extraSpecials, ""))
}

type Page struct {
	Page string `json:"page"`
}

type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\nInput: %s\n", w.Page, strings.Join(w.Words, ", "), w.Input)
}

type Occurrence struct {
	Page  string         `json:"page"`
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	words := []string{}
	for word, occurrence := range o.Words {
		words = append(words, fmt.Sprintf("Word is %s, it showed up %v time", word, occurrence))
	}

	return fmt.Sprintf("Parsed JSON:\nPage: %s\nWords: %s\n", o.Page, strings.Join(words, "\n"))
}

func DoRequest(URL, token string) (Response, error) {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return nil, fmt.Errorf("URL is not valid: %s\nTry add -h flag", err)
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request error: %s", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http Get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %s", err)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("unmarshalling token JSON error: %s", err),
		}
	}

	switch page.Page {
	case "assignment1":
		var parsedAssignment1 Assignment1

		err = json.Unmarshal(body, &parsedAssignment1)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("unmarshalling Assignment JSON error: %s", err),
			}
		}

		return parsedAssignment1, nil

	case "words":
		var parsedWords Words

		err = json.Unmarshal(body, &parsedWords)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("unmarshalling Words JSON error: %s", err),
			}
		}

		return parsedWords, nil

	case "occurrence":
		var parsedOccurrence Occurrence

		err = json.Unmarshal(body, &parsedOccurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("unmarshalling Occurrence JSON error: %s", err),
			}
		}

		return parsedOccurrence, nil
	}

	return nil, nil
}
