package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.Handle("/", logger(http.HandlerFunc(handler)))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain"}
	w.Header()["Cache-Control"] = []string{"no-cache"}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
		return
	}

	tmpl, err := template.New("test").Parse("{{.Message}}\n\nServing from Pod {{.Pod}}")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	resp := response{"Hello Tulsa!", ""}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if "POD_NAME" == pair[0] {
			resp.Pod = pair[1]
		}
	}

	tmpl.Execute(w, resp)
}

type response struct {
	Message string
	Pod     string
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stw := NewStatusTrackingResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(stw, r)
		elapsed := time.Since(start)

		log.Printf(" (%d)  %- 6s %- 10s %s   % 14s\n", stw.StatusCode(), r.Method, r.URL.Path, r.RemoteAddr, elapsed)
	})
}
