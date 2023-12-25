package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var ContentTypeHeaderKey = "Content-Type"
var AppJsonHeaderValue = "application/json"

type Configuration struct {
	Port int
}

type translateRequest struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

func main() {

	config := readConfig("conf.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), mux)
	if err != nil {
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentTypeHeaderKey, AppJsonHeaderValue)
	headerContentType := r.Header.Get(ContentTypeHeaderKey)
	if headerContentType != AppJsonHeaderValue {
		response(w, "Content Type is not "+AppJsonHeaderValue, http.StatusUnsupportedMediaType)
		return
	}
	var t translateRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&t)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			response(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			response(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	result, err := translateTo(t.Text, t.To)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(result)

	w.Write(reqBodyBytes.Bytes())
	return
}

func response(w http.ResponseWriter, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	if len(message) > 0 {
		resp["message"] = message
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
	}
}

func readConfig(fileName string) Configuration {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
