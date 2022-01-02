package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var version string = os.Getenv("VERSION")
var logger = log.New(os.Stdout, "[HttpServer] ", log.LstdFlags | log.Lshortfile)


func main() {
	http.HandleFunc("/", withLogging(index))
	http.HandleFunc("/healthz", withLogging(healthz))
	logger.Print("server listening on 0.0.0.0:80")
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		logger.Fatal("failed to start server")
	}


}

// route handler
func index(w http.ResponseWriter, r *http.Request) {
	header := map[string]string{"X-Server-Version": version}
	for k,v := range r.Header {
		header[k] = strings.Join(v, ", ")
	}

	headerStr, _ := json.Marshal(header)
	var out bytes.Buffer
	json.Indent(&out, headerStr, "", "\t")
	bodyText := "OK"
	responseText := fmt.Sprintf("Body:\n\t%s\nHeader:\n%s\n", bodyText, out.String())
	io.WriteString(w, responseText)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health checks ok"))
}

// HTTPReqInfo describes info about a HTTP request
type HTTPReqInfo struct {
	// GET/POST etc.
	method string
	url string
	referer string
	ipAddr string
	// status code, like 200, 404
	code int
	// number of bytes of the response sent
	size int
	// how long did it take to
	duration time.Duration
	userAgent string
}

type EnhancedResponseWriter struct {
	http.ResponseWriter
	Status int
	ContentLength int
}

func (r *EnhancedResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *EnhancedResponseWriter) Write(data []byte)  (length int, err error) {
	if length, err := r.ResponseWriter.Write(data); err == nil {
		r.ContentLength = length
	}
	return
}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestInfo := &HTTPReqInfo{
			method: r.Method,
			url: r.URL.String(),
			referer: r.Referer(),
			userAgent: r.UserAgent(),
		}
		requestInfo.ipAddr = getRemoteClientAddr(r)
		responseWriter := &EnhancedResponseWriter{
			ResponseWriter: w,
			Status: 200,
		}

		// the actual handler function
		h(responseWriter, r)

		requestInfo.duration = time.Since(start)
		requestInfo.code = responseWriter.Status
		requestInfo.size = responseWriter.ContentLength

		logger.Printf(
			"[ACCESS] %s \"%s %s\" [%d] %dBytes %dms \"%s\" \"%s\"",
			requestInfo.ipAddr,
			requestInfo.method,
			requestInfo.url,
			responseWriter.Status,
			requestInfo.size,
			requestInfo.duration.Milliseconds(),
			requestInfo.referer,
			requestInfo.userAgent,
			)
	}
}

func getRemoteClientAddr(r * http.Request) string {
	// for the case of http server is behind a proxy or LB
	xRealIp := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xRealIp == "" && xForwardedFor == "" {
		// r.RemoteAddress contains port, will be removed
		// i.e: "[::1]:38434" -> "[::1]"
		idx := strings.LastIndex(r.RemoteAddr, ":")
		return r.RemoteAddr[:idx]
	}

	if xForwardedFor != "" {
		xForwardedFors := strings.Split(xForwardedFor, ",")
		// TODO: should filter out non-local address
		// for simplicity, we use the first one for this practice
		return strings.TrimSpace(xForwardedFors[0])
	}

	return xRealIp
}