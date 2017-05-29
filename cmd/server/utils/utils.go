package utils

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/tywkeene/go-agent/version"
)

type APIError struct {
	ErrorMessage string `json:"error_message"`
	HTTPStatus   int    `json:"http_status"`
}

type HttpErrorHandler struct {
	Caller   string
	Response http.ResponseWriter
	Request  *http.Request
}

const (
	ErrorActionErr = iota
	ErrorActionWarn
	ErrorActionDebug
	ErrorActionInfo
)

func NewHttpErrorHandle(caller string, response http.ResponseWriter, request *http.Request) *HttpErrorHandler {
	return &HttpErrorHandler{caller, response, request}
}

//These headers should always be set
func SetResponseHeaders(response http.ResponseWriter, status int) {
	response.Header().Set("Connection", "close")
	response.Header().Set("Server", "Go-Agent server/"+version.GetVersion())
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
}

// HandleError locally, according to the action passed to h.Handle, and then serialized
// in json and sent to the remote address via http, then returns true.
// Otherwise, if there is no error, h.Handle returns false
func (h *HttpErrorHandler) Handle(err error, httpStatus int, action int) bool {
	if err != nil {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		log.Errorf("HttpErrorHandler()->[file:%s line:%d]: %s", file, line, err.Error())
		apiErr := &APIError{
			ErrorMessage: err.Error(),
			HTTPStatus:   httpStatus,
		}
		serialErr, _ := json.Marshal(&apiErr)
		SetResponseHeaders(h.Response, httpStatus)
		io.WriteString(h.Response, string(serialErr))
	}
	return (err != nil)
}

// HandlePanic _Never_ returns on error, instead it panics
func HandlePanic(err error) {
	if err != nil {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		log.Panicf("[file:%s line:%d]: %s", file, line, err.Error())
	}
}

func HandleError(err error, action int) bool {
	if err != nil {
		_, filepath, line, _ := runtime.Caller(1)
		_, file := path.Split(filepath)
		switch action {
		case ErrorActionErr:
			log.Errorf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionWarn:
			log.Warnf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionDebug:
			log.Debugf("[file:%s line:%d]: %s", file, line, err.Error())
			break
		case ErrorActionInfo:
			log.Infof("[file:%s line:%d]: %s", file, line, err.Error())
			break
		}
	}
	return (err != nil)
}

func WriteFile(filename string, source io.Reader) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()
	io.Copy(writer, source)
	return nil
}
