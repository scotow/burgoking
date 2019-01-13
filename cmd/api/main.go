package main

import (
	"github.com/scotow/burgoking"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort                    = 8080
	contentTypeHeader, contentType = "Content-Type", "text/plain"
)

func handle(w http.ResponseWriter, _ *http.Request) {
	code, err := burgoking.GenerateCode(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentTypeHeader, contentType)
	_, _ = w.Write([]byte(code))
}

func listeningAddress() string {
	port, set := os.LookupEnv("PORT")
	if !set {
		port = strconv.Itoa(defaultPort)
	}

	return ":" + port
}

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(listeningAddress(), nil))
}
