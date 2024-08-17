package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Babatunde13/golangr/api/utils"
)

type IHttp interface {
	Get(url string, header interface{}) ([]byte, error)
	Post(url string, header interface{}, d interface{}) ([]byte, error)
	Put(url string, header interface{}, d interface{}) ([]byte, error)
	Delete(url string, header interface{}) ([]byte, error)
}

type Http struct {}

func New () IHttp {
	h := &Http{}
	return IHttp(h)
}

func makeRequest (method string, url string, data interface{}, header interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	client := &http.Client{}
	if data != nil {
		payload, err := json.Marshal(data); if err != nil {

			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload)); if err != nil {

			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil); if err != nil {

			return nil, err
		}
	}

	var headerMap map[string]string = make(map[string]string)
	if header != nil {
		h, ok := header.(map[string]string)
		if !ok {
			err = fmt.Errorf("header is not of type map[string]string")
			fmt.Println("Header is not of type map[string]string")
			return nil, err
		}

		headerMap = h
	}

	contentTypePresent := false
	for key, value := range headerMap {
		req.Header.Set(key, value)
		if strings.ToLower(key) == "content-type" {
			contentTypePresent = true
		}
	}

	if !contentTypePresent {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req); if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body); if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, utils.HandleError(nil, fmt.Sprintf("Error: %s", string(body)))
	}

	return body, nil
}

func (u *Http) Get (url string, header interface{}) ([]byte, error) {
	return makeRequest(http.MethodGet, url, nil, header)
}

func (u *Http) Post (url string, header interface{}, data interface{}) ([]byte, error) {
	return makeRequest(http.MethodPost, url, data, header)
}

func (u *Http) Put (url string, header interface{}, data interface{}) ([]byte, error) {
	return makeRequest(http.MethodPut, url, data, header)
}

func (u *Http) Delete (url string, header interface{}) ([]byte, error) {
	return makeRequest(http.MethodDelete, url, nil, header)
}
