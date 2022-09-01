package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/scotow/burgoking"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

var (
	publicPool  *burgoking.Pool
	privatePool *burgoking.Pool
)

var (
	port    = flag.Int("p", 8080, "listening port")
	contact = flag.String("c", "", "contact address on error")
)

func handleCode(w http.ResponseWriter, r *http.Request) {
	code, _ := burgoking.GenerateCodeStatic(nil)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(code))
	logrus.WithFields(logrus.Fields{"code": code, "ip": realip.FromRequest(r)}).Info("Code used by user.")
}

func handleContact(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(*contact))
}

func main() {
	flag.Parse()

	// Static files.
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/code", handleCode)

	// Contact address.
	if *contact != "" {
		http.HandleFunc("/contact", handleContact)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
