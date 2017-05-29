package connection

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/tywkeene/go-agent/cmd/client/options"
	"github.com/tywkeene/go-agent/version"

	"github.com/tywkeene/go-agent/cmd/server/db"
	"github.com/tywkeene/go-agent/cmd/server/utils"
)

type Connection struct {
	Address        string
	Online         bool
	Authed         bool
	UserAgent      string
	Device         *db.Device
	DeviceRegister *db.DeviceRegister
	client         *http.Client
}

type APIResult struct {
	LocalErr error
	APIErr   *utils.APIError
	Result   []byte
}

func NewAPIResult(localErr error, apiErr *utils.APIError, result []byte) *APIResult {
	return &APIResult{
		LocalErr: localErr,
		APIErr:   apiErr,
		Result:   result,
	}
}

func (r *APIResult) Ok() bool {
	return ((r.LocalErr == nil) && (r.APIErr == nil))
}

func (r *APIResult) PrintErrors() {
	if r.APIErr != nil {
		fmt.Printf("API error: %s (%d)\n", r.APIErr.ErrorMessage, r.APIErr.HTTPStatus)
	}
	if r.LocalErr != nil {
		fmt.Printf("Local error: %s\n", r.LocalErr.Error())
	}
}

func NewConnection(server string) *Connection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return &Connection{
		Address:   server,
		Online:    false,
		Authed:    false,
		UserAgent: "go-agent-" + version.GetVersion(),
		client:    client,
	}
}

func (c *Connection) SetRequestHeaders(request *http.Request) {
	request.Header.Set("Accept-Encoding", "application/x-gzip")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", c.UserAgent)
}

func (c *Connection) ConstructUrl(endpoint string) string {
	urlStr, err := url.Parse(c.Address + endpoint)
	if utils.HandleError(err, utils.ErrorActionErr) == true {
		return ""
	}
	return urlStr.String()
}

func (c *Connection) ConstructPostRequest(endpoint string, data interface{}) *http.Request {
	serial, err := json.Marshal(&data)
	if utils.HandleError(err, utils.ErrorActionErr) == true {
		return nil
	}
	request, err := http.NewRequest("POST", c.ConstructUrl(endpoint), bytes.NewBuffer(serial))
	if utils.HandleError(err, utils.ErrorActionErr) == true {
		return nil
	}
	c.SetRequestHeaders(request)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func inflateResponse(resp *http.Response) ([]byte, error) {
	if resp.Header.Get("Content-Encoding") == "application/x-gzip" {
		reader, err := gzip.NewReader(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		buffer, err := ioutil.ReadAll(reader)
		return buffer, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *Connection) HandleAPIError(response *http.Response, expectStatus int) (*utils.APIError, error) {
	if response.StatusCode != expectStatus {
		buffer, err := inflateResponse(response)
		if err != nil {
			return nil, err
		}
		var errData *utils.APIError
		if err = json.Unmarshal(buffer, &errData); err != nil {
			return nil, err
		}
		return errData, nil
	}
	return nil, nil
}

func (c *Connection) Post(endpoint string, expectStatus int, postData interface{}) *APIResult {
	request := c.ConstructPostRequest(endpoint, postData)
	response, err := c.client.Do(request)
	if err != nil {
		return NewAPIResult(err, nil, nil)
	}
	apiErr, localErr := c.HandleAPIError(response, expectStatus)
	if localErr != nil || apiErr != nil {
		return NewAPIResult(err, apiErr, nil)
	}
	data, err := inflateResponse(response)
	if err != nil {
		return NewAPIResult(err, nil, nil)
	}
	return NewAPIResult(nil, nil, data)
}

func (c *Connection) ConstructDevice(auth *options.Authorization) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	c.Device = &db.Device{
		Hostname: hostname,
		UUID:     auth.UUID,
		AuthStr:  auth.AuthStr,
	}
	return nil
}

func (c *Connection) ConstructDeviceRegister(auth string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	c.DeviceRegister = &db.DeviceRegister{
		Hostname: hostname,
		AuthStr:  auth,
	}
	return nil
}

func (c *Connection) Register(authstr string) *APIResult {
	err := c.ConstructDeviceRegister(authstr)
	if err != nil {
		return NewAPIResult(err, nil, nil)
	}
	result := c.Post("/register", http.StatusOK, &c.DeviceRegister)
	if result.LocalErr != nil {
		return result
	}

	var uuid string
	if err := json.Unmarshal(result.Result, &uuid); err != nil {
		return NewAPIResult(err, nil, nil)
	}

	if uuid == "" {
		return NewAPIResult(fmt.Errorf("Did not get a UUID from the server"), nil, nil)
	}

	auth := &options.Authorization{
		UUID:    uuid,
		AuthStr: authstr,
	}
	c.ConstructDevice(auth)
	return result
}

func (c *Connection) GetStatus(auth *options.Authorization) (bool, error) {
	c.ConstructDevice(auth)
	// We need to catch the http status code so we do this semi manually
	request := c.ConstructPostRequest("/status", c.Device)
	response, err := c.client.Do(request)
	if response.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}

func (c *Connection) Login() *APIResult {
	return c.Post("/login", http.StatusOK, &c.Device)
}

func (c *Connection) Logout() *APIResult {
	return c.Post("/logout", http.StatusOK, &c.Device)
}

func (c *Connection) Ping() *APIResult {
	return c.Post("/ping", http.StatusOK, &c.Device)
}
