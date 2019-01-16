package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort                    = 8080
	contentTypeHeader, contentType = "Content-Type", "text/plain"
)

var (
	pool = NewPool()
)

func handleCode(w http.ResponseWriter, r *http.Request) {
	codeC, cancelC := make(chan string), make(chan struct{})
	go pool.GetCode(codeC, cancelC)

	select {
	case code := <-codeC:
		w.WriteHeader(http.StatusOK)
		w.Header().Set(contentTypeHeader, contentType)
		_, _ = w.Write([]byte(code))
		logrus.WithFields(logrus.Fields{"code": code, "ip": r.RemoteAddr}).Info("Code used by user.")
	case <-r.Context().Done():
		cancelC <- struct{}{}
	}

}

func handleDirect(w http.ResponseWriter, r *http.Request) {

}

func listeningAddress() string {
	port, set := os.LookupEnv("PORT")
	if !set {
		port = strconv.Itoa(defaultPort)
	}

	return ":" + port
}

func main() {
	go pool.fill()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/code", handleCode)
	http.HandleFunc("/direct", handleDirect)

	log.Fatal(http.ListenAndServe(listeningAddress(), nil))
}
