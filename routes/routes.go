package routes

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tywkeene/go-tracker/auth"
	"github.com/tywkeene/go-tracker/db"
	"github.com/tywkeene/go-tracker/utils"

	"github.com/satori/go.uuid"
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
	response.Header().Set("Connection", "close")
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
	errHandle := utils.NewHttpErrorHandle("registerHandle", w, r)
	if validateRequestMethod(errHandle, "POST") == false {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var registerAuth *db.DeviceRegister
	err := decoder.Decode(&registerAuth)
	if errHandle.Handle(err, http.StatusInternalServerError, utils.ErrorActionErr) == true {
		return
	}

	err = auth.ValidateRegisterAuth(registerAuth.AuthStr)
	if errHandle.Handle(err, http.StatusUnauthorized, utils.ErrorActionErr) == true {
		return
	}
	deviceUUID := uuid.NewV4().String()
	uuidJson, err := json.Marshal(deviceUUID)
	if errHandle.Handle(err, http.StatusInternalServerError, utils.ErrorActionErr) == true {
		return
	}

	addr := strings.Split(r.RemoteAddr, ":")
	device := &db.Device{
		UUID:        deviceUUID,
		Address:     addr[0],
		AuthStr:     registerAuth.AuthStr,
		Hostname:    registerAuth.Hostname,
		Online:      false,
		LastSeen:    nil,
		LocationLog: nil,
	}

	err = db.HandleRegister(device)
	if errHandle.Handle(err, http.StatusUnauthorized, utils.ErrorActionErr) == true {
		return
	}

	log.Printf("Hostname '%s' successfully registered (authstr:%s) (uuid:%s)",
		device.Hostname, registerAuth.AuthStr, device.UUID)

	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	setDefaultResponseHeaders(w)
	io.WriteString(w, string(uuidJson))
}

func pingHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func loginHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
	errHandle := utils.NewHttpErrorHandle("loginHandle", w, r)
	if validateRequestMethod(errHandle, "POST") == false {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var device *db.Device
	err := decoder.Decode(&device)
	if errHandle.Handle(err, http.StatusInternalServerError, utils.ErrorActionErr) == true {
		return
	}

	online, err := db.IsDeviceOnline(device)
	if err != nil {
		log.Println(err)
		errHandle.Handle(fmt.Errorf("error getting device status"),
			http.StatusInternalServerError, utils.ErrorActionErr)
		return
	}
	if online == true {
		errHandle.Handle(fmt.Errorf("that device is already online"), http.StatusBadRequest, utils.ErrorActionErr)
		return
	}

	err = db.HandleLogin(device)
	if errHandle.Handle(err, http.StatusUnauthorized, utils.ErrorActionErr) == true {
		return
	}
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	setDefaultResponseHeaders(w)
	log.Printf("Device '%s' logged in [authstr:%s] [uuid:%s]",
		device.Hostname, device.AuthStr, device.UUID)
}

func logoffHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
	errHandle := utils.NewHttpErrorHandle("logoffHandle", w, r)
	if validateRequestMethod(errHandle, "POST") == false {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var device *db.Device
	err := decoder.Decode(&device)
	if errHandle.Handle(err, http.StatusInternalServerError, utils.ErrorActionErr) == true {
		return
	}

	online, err := db.IsDeviceOnline(device)
	if err != nil {
		log.Println(err)
		errHandle.Handle(fmt.Errorf("error getting device status"),
			http.StatusInternalServerError, utils.ErrorActionErr)
		return
	}
	if online == false {
		errHandle.Handle(fmt.Errorf("that device is already offline"), http.StatusBadRequest, utils.ErrorActionErr)
		return
	}

	err = db.HandleLogoff(device)
	if errHandle.Handle(err, http.StatusUnauthorized, utils.ErrorActionErr) == true {
		return
	}
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	setDefaultResponseHeaders(w)

	log.Printf("Device '%s' logged off [authstr:%s] [uuid:%s]",
		device.Hostname, device.AuthStr, device.UUID)
}

func errorHandle(w http.ResponseWriter, r *http.Request) {
	LogHttp(r)
}

func RegisterHandles() {
	http.HandleFunc("/register", registerHandle)
	http.HandleFunc("/ping", pingHandle)
	http.HandleFunc("/login", loginHandle)
	http.HandleFunc("/logoff", logoffHandle)
	http.HandleFunc("/report_error", errorHandle)
}

func Launch() {
	panic(http.ListenAndServe(":8080", nil))
}
