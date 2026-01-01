package server

import (
	// "fmt"
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
}

type U64Server struct {
	creds U64Creds
}

type U64ResponseBase struct {
	Errors []string `json:"errors"`
}

// type U64Version struct {
// 	Errors []string `json:"errors"`
// }

func NewU64Server(creds U64Creds) (U64Server, error) {
	server := U64Server{}
	server.creds = creds

	// test server connection
	_, err := server.RestCall("GET", "/v1/version", nil)
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

	resp, err := self.RestCallRaw(req)
	if err != nil { return err }

	_ = resp

	return nil
}

func (self *U64Server) RestCall(method, path string, params map[string]string) ([]byte, error) {
	req, err := self.CreateRestRequest(method, path, nil)
	if err != nil { return nil, err }

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	data, err := self.RestCallRaw2(req)
	if err != nil { return nil, err }

	// err = json.NewDecoder(body).Decode(&result)
	// if err != nil { return result, err }

	return data, nil
}

func (self *U64Server) PeekMemory(address uint16, length uint32) ([]byte, error) {
	data, err := self.RestCall("GET", "/v1/machine:readmem", map[string]string{
		"address": fmt.Sprintf("%04x", address),
		"length":  fmt.Sprintf("%d", length),
	})

	if err != nil { return data, err }

	return data, nil
}

func (self *U64Server) PokeMemory(address uint16, data []byte) error {
	hex_string := ""
	for _, b := range data {
		hex_string += fmt.Sprintf("%02x", b)
	}

	fmt.Println(hex_string)

	data, err := self.RestCall("PUT", "/v1/machine:writemem", map[string]string{
		"address": fmt.Sprintf("%04x", address),
		"data": hex_string,
	})

	if err != nil { return err }

	return nil
}

func (self *U64Server) RestCallRaw2(req *http.Request) ([]byte, error) {
	// var result any

	resp, err := self.RestCallRaw(req)
	if err != nil { return nil, err }

	// length := resp.ContentLength
	data := make([]byte, resp.ContentLength)

	_, err = io.ReadAtLeast(resp.Body, data, int(resp.ContentLength))
	if err != nil { return nil, err }

	// fmt.Println(result.())

	return data, nil
}

func (self *U64Server) RestCallRaw(req *http.Request) (*http.Response, error) {
	// var result any
	resp, err := http.DefaultClient.Do(req)

	if err != nil { return nil, err }

	return resp, nil
}

