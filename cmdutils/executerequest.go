package cmdutil

import (
	"github.com/go-resty/resty/v2"
	"log"
)

func ExecutePostRequest(url string, body []byte) (*resty.Response, error) {
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(url)
	return resp, err
}

func HandleResponse(resp *resty.Response, expectedCode int) {
	if resp.StatusCode() != expectedCode {
		log.Printf("некорректный статус код: %s\n", resp.Status())
	}
	log.Println(resp.String())
}
