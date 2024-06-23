package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Response interface {
	GetResponse() string
}

type RequestDetails struct {
	Token    string
	Password string
	URL      string
	Client   http.Client
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

	response, err := r.Client.Post(parsedURL.Scheme+"://"+parsedURL.Host+"/login", "application/json", bytes.NewBuffer(jsonBody))
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

	fmt.Printf("GetToken token request response body is: %v\n", string(body))

	var retrievedToken Token

	err = json.Unmarshal(body, &retrievedToken)
	if err != nil {
		return RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("unmarshalling token JSON error: %s", err),
		}
	}

	fmt.Printf("GetToken retrievedToken.Token is: %v\n", retrievedToken.Token)

	r.Token = retrievedToken.Token

	return nil
}

type Token struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

func DoLogin(requestDetails RequestDetails) (Response, error) {
	if requestDetails.Token == "" {
		err := requestDetails.GetToken()
		fmt.Printf("DoLogin requestDetails.Token: %v\n", requestDetails.Token)
		if err != nil {
			return nil, err
		}
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
			err := requestDetails.GetToken()
			if err != nil {
				return nil, err
			}
		} else {
			fmt.Println("Token is still valid, re-use")
		}
	}

	return requestDetails, nil
}
