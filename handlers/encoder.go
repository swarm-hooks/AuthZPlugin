package handlers

import (
	"encoding/json"
	"net/http"
)

// DefaultContentTypeV1_1 is the default content type accepted and sent by the plugins.
const DefaultContentTypeV1_1 = "application/vnd.docker.plugins.v1.1+json"

// DecodeRequest decodes an http request into a given structure.
func DecodeRequest(w http.ResponseWriter, r *http.Request, req interface{}) (err error) {
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

// EncodeResponse encodes the given structure into an http response.
func EncodeResponse(w http.ResponseWriter, res interface{}, err string) {
	w.Header().Set("Content-Type", DefaultContentTypeV1_1)
	if err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}
