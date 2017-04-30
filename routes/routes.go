package routes

import (
	"fmt"
	"log"
	"net/url"
	//	"encoding/json"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	//	"github.com/tywkeene/go-tracker/db"
	"github.com/tywkeene/go-tracker/utils"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//Handle and make sure the client wants or can handle gzip, and replace the writer if it
//can, if not, simply use the normal http.ResponseWriter
func GzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "application/x-gzip") == false {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "application/x-gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func LogHttp(r *http.Request) {
	log.Printf("%s %s %s %s", r.Method, r.URL, r.RemoteAddr, r.UserAgent())
}

//These headers should always be set
func setDefaultResponseHeaders(response http.ResponseWriter) {
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("Server", "Go Tracker v0.0.0")
}

//Checks a request header and ensures it is allowed, otherwise it will set the Allow http header
// and return HTTP 405 Method Not Allowed
func validateRequestMethod(errHandle *utils.HttpErrorHandler, allowed string) bool {
	if strings.Contains(allowed, errHandle.Request.Method) == false {
		errHandle.Response.Header().Set("Allow", allowed)
		setDefaultResponseHeaders(errHandle.Response)
		errHandle.Handle(fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed, utils.ErrorActionErr)
		return false
	}
	return true
}

//GetQueryValue() takes a name of a key:value pair to fetch from a URL encoded query,
//a http.ResponseWriter 'w', and a http.Request 'r'. In the event that an error is encountered
//the error will be returned to the client via logging facilities that use 'w' and 'r'
func GetQueryValue(name string, w http.ResponseWriter, r *http.Request) (string, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if query == nil || err != nil {
		return "", err
	}
	return query.Get(name), nil
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func pingHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func loginHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func logoffHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func errorHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func RegisterHandles() {
	http.HandleFunc("/register", registerHandle)
	http.HandleFunc("/ping", pingHandle)
	http.HandleFunc("/register_login", loginHandle)
	http.HandleFunc("/register_logoff", logoffHandle)
	http.HandleFunc("/register_error", errorHandle)
}

func Launch() {
	panic(http.ListenAndServe(":8080", nil))
}
