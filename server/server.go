package server

import (
	// "fmt"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	// "io"
	"net/http"
	"net/url"
	// "os"
	// "encoding/json"
)

type U64Creds struct {
	Address, Password string
	EnableMessageBox bool
}

type U64Server struct {
	creds U64Creds
}

type U64ResponseBase struct {
	Errors []string `json:"errors"`
}

func NewU64Server(creds U64Creds) (U64Server, error) {
	server := U64Server{}
	server.creds = creds

	// test server connection
	_, err := server.RestCall("GET", "/v1/version", nil, nil)
	if err != nil { return server, err }

	return server, nil
}

func (self *U64Server) CreateRestRequest(method, path string, body io.Reader) (*http.Request, error) {
	path, err := url.JoinPath(self.creds.Address, path)
	if err != nil { return nil, err }

	req, err := http.NewRequest(method, path, body)
	if err != nil { return nil, err }

	req.Header.Add("X-Password", self.creds.Password)

	return req, nil
}

func (self *U64Server) RunPRG(reader io.Reader) error {
	req, err := self.CreateRestRequest("POST", "/v1/runners:run_prg", reader)
	if err != nil { return err }

	req.Header.Set("Content-Type", "application/octet-stream")

	_, err = self.RestCallRaw(req)
	if err != nil { return err }


	return nil
}

func (self *U64Server) RestCall(method, path string, reader io.Reader, params map[string]string) ([]byte, error) {
	req, err := self.CreateRestRequest(method, path, reader)
	if err != nil { return nil, err }

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	data, err := self.RestCallRaw(req)
	if err != nil { return nil, err }

	// err = json.NewDecoder(body).Decode(&result)
	// if err != nil { return result, err }

	return data, nil
}

func (self *U64Server) PeekMemory(address uint16, length uint32) ([]byte, error) {
	data, err := self.RestCall("GET", "/v1/machine:readmem", nil, map[string]string{
		"address": fmt.Sprintf("%04x", address),
		"length":  fmt.Sprintf("%d", length),
	})

	if err != nil { return data, err }

	return data, nil
}

func (self *U64Server) PokeMemoryWithStream(address uint16, reader io.Reader) error {
	data, err := self.RestCall("POST", "/v1/machine:writemem", reader, map[string]string{
		"address": fmt.Sprintf("%04x", address),
	})

	_ = data

	if err != nil { return err }

	return nil
}

func (self *U64Server) PokeMemory(address uint16, data []byte) error {
	return self.PokeMemoryWithStream(address, bytes.NewBuffer(data))
}

func (self *U64Server) RestCallRaw(req *http.Request) ([]byte, error) {
	// var result any

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }

	// length := resp.ContentLength
	data := make([]byte, resp.ContentLength)

	_, err = io.ReadAtLeast(resp.Body, data, int(resp.ContentLength))
	if err != nil { return nil, err }

	if resp.Header.Get("Content-Type") == "application/json" {
		rest_response := U64ResponseBase{}
		err = json.Unmarshal(data, &rest_response)
		if err != nil { return nil, err }

		if len(rest_response.Errors) > 0 {
			err = nil

			for i, err_text := range rest_response.Errors {
				new_err := fmt.Errorf("REST Error %d: %s", i, err_text)
				if err == nil {
					err = new_err
				} else {
					err = errors.Join(err, new_err)
				}
			}
			return nil, err
		}
	}

	// fmt.Println(result.())

	return data, nil
}

