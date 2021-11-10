package api

import (
	"encoding/json"
	"magnusquiz/pkg/log"
	"net/http"
)

const badRequest = "Bad request"
const iseRequest = "Internal server error"
const unauthorized = "unauthorized"

var Key string

type HTTPResponse struct {
	err  string
	data interface{}
	w    http.ResponseWriter
}

func (r *HTTPResponse) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data":  r.data,
		"error": r.err,
	})
}

func (r *HTTPResponse) checkQkey(re *http.Request) bool {
	t := re.URL.Query().Get("k")
	return t == Key
}

func (r *HTTPResponse) unauthorized() {
	r.err = unauthorized
	r.w.WriteHeader(http.StatusUnauthorized)
	r.write()
}

func (r *HTTPResponse) ok(d interface{}) {
	r.data = d
	r.w.WriteHeader(http.StatusOK)
	r.write()
}

func (r *HTTPResponse) badRequest() {
	r.err = badRequest
	r.w.WriteHeader(http.StatusBadRequest)
	r.write()
}

func (r *HTTPResponse) ise(err error) {
	log.Logger.Errorf("error on api. ERR => %s \n", err)
	r.err = iseRequest
	r.w.WriteHeader(http.StatusInternalServerError)
	r.write()
}

func (r *HTTPResponse) write() {
	b, err := r.Marshal()
	if err != nil {
		log.Logger.Errorf("error in marshaling results. ERR => %s \n", err)
		_, _ = r.w.Write([]byte{})
		return
	}
	_, err = r.w.Write(b)
	if err != nil {
		log.Logger.Errorf("error in writing results to http.ResponseWriter. ERR => %s \n", err)
	}
}
