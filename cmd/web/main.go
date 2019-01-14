package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort                    = 8080
	contentTypeHeader, contentType = "Content-Type", "text/plain"
)

func handle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("index.html"))
}

func handleCode(w http.ResponseWriter, r *http.Request) {

	select {
	case <-r.Context().Done():
		fmt.Println("Request canceled.")
	case <-time.After(5 * time.Second):
		fmt.Println("Request succeeded.")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentTypeHeader, contentType)
	w.Write([]byte("Hello, World!"))
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
	http.HandleFunc("/code", handleCode)
	log.Fatal(http.ListenAndServe(listeningAddress(), nil))
}
