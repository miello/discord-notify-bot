package utils

import (
	"api-gateway/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type HttpType string

const (
	POST   HttpType = "POST"
	GET    HttpType = "GET"
	PUT    HttpType = "PUT"
	DELETE HttpType = "DELETE"
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

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, fmt.Errorf("error with status code: %v", res.StatusCode)
	}

	// doc, err := goquery.NewDocumentFromReader(res.Body)

	// el := doc.Find("*:contains('Please login with either of the following choices.')")

	// if el.Length() > 0 {
	// 	log.Fatalf("Please login with either of the following choices.")
	// 	return nil, fmt.Errorf("Please contact to bot owner to change key")
	// }

	return res, err
}

func GetJSONByFormDataReq(httpType HttpType, path string, req_body *map[string]string, res_body interface{}) error {
	BASE_URL := os.Getenv("BASE_URL")
	SESSION_KEY := os.Getenv("SESSION_KEY")
	SESSION_VALUE := os.Getenv("SESSION_VALUE")

	session_cookie := &http.Cookie{
		Name:  SESSION_KEY,
		Value: SESSION_VALUE,
	}

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	for key, value := range *req_body {
		fw, _ := writer.CreateFormField(key)
		io.Copy(fw, strings.NewReader(value))
	}

	req, err := http.NewRequest(string(httpType), fmt.Sprintf("%v/%v", BASE_URL, path), bytes.NewReader(form.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err != nil {
		return err
	}

	req.AddCookie(session_cookie)
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return fmt.Errorf("error with status code: %v", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(res_body)
}

func ExtractError(err error) (int, types.ResponseError) {
	arr := strings.Split(err.Error(), ": ")
	status_code, _ := strconv.Atoi(arr[0])

	msg := strings.Join(arr[1:], " ")

	return status_code, types.ResponseError{
		Msg: msg,
	}
}

func CreateError(status_code int, msg string) error {
	return fmt.Errorf("%v: %v", status_code, msg)
}
