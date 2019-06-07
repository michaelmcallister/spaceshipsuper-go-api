package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DoGet will send a HTTP GET request to the supplied endpoint. It requires the
// current bearer Oauth2 token.
func DoGet(endpoint, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 error codes.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recieved HTTP %d from %s", resp.StatusCode, endpoint)
	}

	return ioutil.ReadAll(resp.Body)
}

// DoPost will send a HTTP POST request to the supplied endpoint, all known
// endpoints for the API require a token which is a base64 representation of the
// username and password. And optionally a byte array of postData.
func DoPost(endpoint, token string, postData []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 error codes.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recieved HTTP %d from %s", resp.StatusCode, endpoint)
	}

	return ioutil.ReadAll(resp.Body)
}
