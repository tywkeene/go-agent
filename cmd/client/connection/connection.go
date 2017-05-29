package connection

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tywkeene/go-tracker/cmd/client/options"
	"github.com/tywkeene/go-tracker/version"

	"github.com/tywkeene/go-tracker/cmd/server/db"
	"github.com/tywkeene/go-tracker/cmd/server/utils"
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

func NewConnection(server string) *Connection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return &Connection{
		Address:   server,
		Online:    false,
		Authed:    false,
		UserAgent: "go-tracker-" + version.GetVersion(),
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
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		buffer, err := ioutil.ReadAll(reader)
		return buffer, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *Connection) HandleAPIError(response *http.Response, expectStatus int) error {
	if response.StatusCode != expectStatus {
		defer response.Body.Close()
		buffer, err := inflateResponse(response)
		if err != nil {
			return err
		}
		var errData *utils.APIError
		if err = json.Unmarshal(buffer, &errData); err != nil {
			return err
		}
		return fmt.Errorf("Error [%s]->(HTTP %d %s): %s",
			c.Address, errData.HTTPStatus, http.StatusText(errData.HTTPStatus), errData.ErrorMessage)
	}
	return nil
}

func (c *Connection) Post(endpoint string, expectStatus int, data interface{}) ([]byte, error) {
	request := c.ConstructPostRequest(endpoint, data)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if err := c.HandleAPIError(response, expectStatus); err != nil {
		return nil, err
	}
	return inflateResponse(response)
}

func (c *Connection) constructDevice(auth *options.Authorization) error {
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

func (c *Connection) constructDeviceRegister(auth string) error {
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

func (c *Connection) Register(authstr string) error {
	err := c.constructDeviceRegister(authstr)
	if err != nil {
		return err
	}
	response, err := c.Post("/register", http.StatusOK, &c.DeviceRegister)
	if err != nil {
		return err
	}

	var uuid string
	if err := json.Unmarshal(response, &uuid); err != nil {
		return err
	}

	if uuid == "" {
		return fmt.Errorf("Did not get a UUID from the server")
	}

	auth := &options.Authorization{
		UUID:    uuid,
		AuthStr: authstr,
	}
	c.constructDevice(auth)
	log.Println("Successfully registered with", c.Address)
	return nil
}

func (c *Connection) GetStatus(auth *options.Authorization) (bool, error) {
	c.constructDevice(auth)
	// We need to catch the http status code so we do this semi manually
	request := c.ConstructPostRequest("/status", c.Device)
	response, err := c.client.Do(request)
	if response.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}

func (c *Connection) Login() error {
	_, err := c.Post("/login", http.StatusOK, &c.Device)
	return err
}

func (c *Connection) Logout() error {
	_, err := c.Post("/logout", http.StatusOK, &c.Device)
	return err
}

func (c *Connection) Ping() error {
	_, err := c.Post("/ping", http.StatusOK, &c.Device)
	return err
}
