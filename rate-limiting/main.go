package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	// construct default logger
	var programLevel = new(slog.LevelVar) // Info by default
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(logger))

	// set log level to debug, if OS env DEBUG set as 1
	if os.Getenv("DEBUG") == "1" {
		programLevel.Set(slog.LevelDebug)
	}

	var requestCounter int

	// 5 ticks per second
	ticker := time.NewTicker(time.Second / 5)
	defer ticker.Stop()
	done := make(chan bool)

	// during 10 sec
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			sentRequest(&requestCounter)
		}
	}
}

func sentRequest(requestCounter *int) {
	slog.Info("requestCounter is " + fmt.Sprintf("%v", *requestCounter))
	response, err := http.Get("http://localhost:8080/ratelimit")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info(string(responseBody))

	*requestCounter++
}
