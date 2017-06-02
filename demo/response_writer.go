package main

import "net/http"

type StatusTrackingResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func NewStatusTrackingResponseWriter(w http.ResponseWriter) *StatusTrackingResponseWriter {
	return &StatusTrackingResponseWriter{ResponseWriter: w}
}

func (w *StatusTrackingResponseWriter) StatusCode() int {
	return w.status
}

func (w *StatusTrackingResponseWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *StatusTrackingResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}
