package api

import "fmt"

type RequestError struct {
	HTTPCode int
	Body     string
	Err      string
}

func (r RequestError) Error() string {
	return fmt.Sprintf("%v\nResponse code: %v\nResponse body: %v\n", r.Err, r.HTTPCode, r.Body)
}
