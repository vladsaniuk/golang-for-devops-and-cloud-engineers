package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/vladsaniuk/golang-for-devops-and-cloud-engineers/http-login-request/pkg/api"
)

func main() {
	var (
		URL      string
		password string
		token    string
		count    int
	)

	// construct default logger
	var programLevel = new(slog.LevelVar) // Info by default
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(logger))

	// set log level to debug, if OS env DEBUG set as 1
	if os.Getenv("DEBUG") == "1" {
		programLevel.Set(slog.LevelDebug)
	}

	flag.StringVar(&URL, "url", "", "URL to fetch")
	flag.StringVar(&password, "password", "", "Password to use to get token for the API calls")
	flag.IntVar(&count, "count", 1, "number of request")
	flag.Parse()

	if password == "" {
		slog.Error("Please, provide password\nTry add -h flag")
		os.Exit(1)
	}

	requestDetails := api.RequestDetails{
		Token:    token,
		Password: password,
		URL:      URL,
		Client:   http.Client{},
	}

	sum := 0
	for i := 1; i <= count; i++ {
		getToken, err := api.DoLogin(requestDetails)
		if err != nil {
			slog.Error("error making login request: " + err.Error())
			os.Exit(1)
		}

		requestDetails.Token = getToken.GetResponse()

		slog.Debug("main requestDetails.Token: " + requestDetails.Token)

		response, err := api.DoRequest(requestDetails)
		slog.Debug("main response: " + response.GetResponse())
		if err != nil {
			slog.Error("error making request: " + err.Error())
			os.Exit(1)
		} else if response == nil {
			slog.Error("Something went wrong - got nil in response")
			os.Exit(1)
		}

		fmt.Println(response.GetResponse())

		sum += i
	}
}
