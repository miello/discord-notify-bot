package utils

import (
	"fmt"
	"net/http"
	"os"
)

var client http.Client

func GetHTML(path string) (*http.Response, error) {
	BASE_URL := os.Getenv("BASE_URL")
	SESSION_KEY := os.Getenv("SESSION_KEY")
	SESSION_VALUE := os.Getenv("SESSION_VALUE")

	session_cookie := &http.Cookie{
		Name:  SESSION_KEY,
		Value: SESSION_VALUE,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/%v", BASE_URL, path), nil)

	if err != nil {
		return nil, err
	}

	req.AddCookie(session_cookie)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, err
}
