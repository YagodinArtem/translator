package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type translateRequest struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(":8987", mux)
	if err != nil {
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		response(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
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
	//w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	if len(message) > 0 {
		resp["message"] = message
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
	}
}
