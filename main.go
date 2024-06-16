package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

type Response interface {
	GetResponse() string
}

type RequestDetails struct {
	Token    string
	Password string
	URL      string
}

func (r RequestDetails) GetResponse() string {
	return fmt.Sprintf("%v", r.Token)
}

func (r *RequestDetails) GetToken() error {
	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil {
		fmt.Printf("error parsing request URL: %s\n", err)
		os.Exit(1)
	}

	loginRequest := LoginRequest{
		Password: r.Password,
	}

	jsonBody, err := json.Marshal(loginRequest)
	if err != nil {
		return fmt.Errorf("error marshaling password to JSON: %s", err)
	}

	response, err := http.Post(parsedURL.Scheme+"://"+parsedURL.Host+"/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error making POST request: %s", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error making POST request, statue code is %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %s", err)
	}

	var retrievedToken Token

	err = json.Unmarshal(body, &retrievedToken)
	if err != nil {
		return fmt.Errorf("unmarshalling token JSON error: %s", err)
	}

	r.Token = retrievedToken.Token

	return nil
}

type Token struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Password string `json:"password"`
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

func main() {
	var (
		URL      string
		password string
		token    string
		count    int
	)

	flag.StringVar(&URL, "url", "", "URL to fetch")
	flag.StringVar(&password, "password", "", "Password to use to get token for the API calls")
	flag.IntVar(&count, "count", 1, "number of request")
	flag.Parse()

	if password == "" {
		fmt.Println("Please, provide password\nTry add -h flag")
		os.Exit(1)
	}

	sum := 0
	for i := 1; i <= count; i++ {
		getToken, err := doLogin(URL, token, password)
		if err != nil {
			fmt.Printf("error making login request: %s\n", err)
			os.Exit(1)
		}

		token = getToken.GetResponse()

		response, err := doRequest(URL, token)
		if err != nil {
			fmt.Printf("error making request: %s\n", err)
			os.Exit(1)
		} else if response == nil {
			fmt.Println("Something went wrong - got nil in response")
			os.Exit(1)
		}

		fmt.Printf("%s\n", response.GetResponse())

		sum += i
	}
}

func doRequest(URL, token string) (Response, error) {
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
		return nil, fmt.Errorf("unmarshalling Page JSON error: %s\nGot the following response from server: %s", err, string(body))
	}

	switch page.Page {
	case "assignment1":
		var parsedAssignment1 Assignment1

		err = json.Unmarshal(body, &parsedAssignment1)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling Assignment JSON error: %s", err)
		}

		return parsedAssignment1, nil

	case "words":
		var parsedWords Words

		err = json.Unmarshal(body, &parsedWords)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling Words JSON error: %s", err)
		}

		return parsedWords, nil

	case "occurrence":
		var parsedOccurrence Occurrence

		err = json.Unmarshal(body, &parsedOccurrence)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling Occurrence JSON error: %s", err)
		}

		return parsedOccurrence, nil
	}

	return nil, nil
}

func doLogin(URL, token, password string) (Response, error) {
	requestDetails := RequestDetails{
		Token:    token,
		Password: password,
		URL:      URL,
	}

	if requestDetails.Token == "" {
		requestDetails.GetToken()
	} else {
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(requestDetails.Token, jwt.MapClaims{})
		if err != nil {
			return nil, err
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("parsedToken.Claims is not of type jwt.MapClaims")
		}

		expirationFloat64, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("claims is not of type float64")
		}

		expiration := time.Unix(int64(expirationFloat64), 0)
		if time.Now().After(expiration) {
			fmt.Println("Token expired, do refresh")
			requestDetails.GetToken()
		} else {
			fmt.Println("Token is still valid, re-use")
		}
	}

	return requestDetails, nil
}
