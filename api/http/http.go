package http

import (
	"bkoiki950/go-store/api/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func makeRequest (method string, url string, data interface{}) (interface{}, error) {
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

	req.Header.Set("Content-Type", "application/json")
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

	var result interface{}
	json.Unmarshal(body, &result)
	return result, nil
}

func Get (url string) (data interface{}, err error) {
	return makeRequest(http.MethodGet, url, nil)
}

func Post (url string, data interface{}) (interface{}, error) {
	return makeRequest(http.MethodPost, url, data)
}

func Put (url string, data interface{}) (interface{}, error) {
	return makeRequest(http.MethodPut, url, data)
}

func Delete (url string) (interface{}, error) {
	return makeRequest(http.MethodDelete, url, nil)
}
