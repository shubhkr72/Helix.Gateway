package middleware

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	Status int
	Bytes  int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.Status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.Bytes += n
	return n, err
}