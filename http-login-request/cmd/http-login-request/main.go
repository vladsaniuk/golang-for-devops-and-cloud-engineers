package main

import (
	"flag"
	"fmt"
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
		getToken, err := api.DoLogin(URL, token, password)
		if err != nil {
			fmt.Printf("error making login request: %s\n", err)
			os.Exit(1)
		}

		token = getToken.GetResponse()

		response, err := api.DoRequest(URL, token)
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
