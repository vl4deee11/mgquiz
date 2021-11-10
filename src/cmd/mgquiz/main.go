package main

import (
	"fmt"
	"magnusquiz/internal/api"
	"magnusquiz/pkg/log"
	"magnusquiz/pkg/repo"
	"net/http"
	"regexp"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port uint16
	Key  string
	DB   struct {
		Port uint16
		Host string
		User string
		Pass string
		Name string
	} `split_words:"true"`
	LogLVL string `default:"INFO"`
}

func main() {
	c := &Config{}
	err := envconfig.Process("MG", c)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.MakeLogger(c.LogLVL)
	repo.Connect(c.DB.Host, c.DB.User, c.DB.Pass, c.DB.Name, c.DB.Port)
	api.Key = c.Key
	mux := http.NewServeMux()
	mux.HandleFunc("/", route)

	// Add the pprof routes
	//mux.HandleFunc("/debug/pprof/", pprof.Index)
	//mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	//
	//mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	//mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	//mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	//mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	log.Logger.Info("service is up")
	log.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), mux))
}

type WriterWithStatusCode struct {
	http.ResponseWriter
	StatusCode int
}

func NewWriter(w http.ResponseWriter) *WriterWithStatusCode {
	return &WriterWithStatusCode{w, http.StatusOK}
}

func (w *WriterWithStatusCode) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

var rGenQuestions = regexp.MustCompile(`/api/questions/generate.*`)
var rQuestions = regexp.MustCompile(`/api/questions`)
var rUserInfo = regexp.MustCompile(`/api/user-info.*`)
var rAnswers = regexp.MustCompile(`/api/answers`)

func route(w http.ResponseWriter, r *http.Request) {
	lw := NewWriter(w)
	switch {
	case r.URL.Path == "/health" || r.URL.Path == "/ready":
		lw.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(lw, "ok")
	case rGenQuestions.MatchString(r.URL.Path):
		if r.Method == http.MethodGet {
			api.GenerateQuestions(lw, r)
		} else {
			lw.WriteHeader(http.StatusMethodNotAllowed)
		}
	case rQuestions.MatchString(r.URL.Path):
		if r.Method == http.MethodPost {
			api.CreateQuestion(lw, r)
		} else {
			lw.WriteHeader(http.StatusMethodNotAllowed)
		}
	case rAnswers.MatchString(r.URL.Path):
		if r.Method == http.MethodPost {
			api.CreateAnswer(lw, r)
		} else {
			lw.WriteHeader(http.StatusMethodNotAllowed)
		}
	case rUserInfo.MatchString(r.URL.Path):
		if r.Method == http.MethodPost {
			api.CreateUserInfo(lw, r)
		} else if r.Method == http.MethodGet {
			api.GetUserInfo(lw, r)
		} else {
			lw.WriteHeader(http.StatusMethodNotAllowed)
		}
	default:
		lw.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(lw, "404 not found")
	}
	log.Logger.Infof(
		"%s %s%s [%d]",
		r.Method,
		r.Host,
		r.URL.String(),
		lw.StatusCode,
	)
}
