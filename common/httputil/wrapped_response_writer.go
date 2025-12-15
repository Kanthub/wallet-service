package httputil

import (
	"encoding/json"
	"net/http"
)

type WrappedResponseWriter struct {
	StatusCode  int
	ResponseLen int

	w           http.ResponseWriter
	wroteHeader bool
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{
		StatusCode: 200,
		w:          w,
	}
}

func (w *WrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *WrappedResponseWriter) Write(bytes []byte) (int, error) {
	n, err := w.w.Write(bytes)
	w.ResponseLen += n
	return n, err
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.StatusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func WriteSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": data,
	})
}

func WriteError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]any{
		"code": -1,
		"msg":  msg,
	})
}
