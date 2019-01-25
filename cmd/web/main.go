package main

import (
	"flag"
	"github.com/scotow/burgoking"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	publicPool 	*burgoking.Pool
	privatePool 	*burgoking.Pool
)

var (
	port 			= flag.Int("p", 8080, "listening port")
	contact 		= flag.String("c", "", "contact address on error")

	publicSize 		= flag.Int("n", 3, "public code pool size")
	publicExpiration 	= flag.String("d", time.Duration(24 * time.Hour).String(), "public code expiration")
	publicRetry 		= flag.String("r", time.Duration(30 * time.Second).String(), "public code regeneration interval")

	privateDirectKey	= flag.String("k", "", "authorization token for private and direct code (disable if empty)")
	privateSize 		= flag.Int("N", 1, "private code pool size")
	privateExpiration 	= flag.String("D", time.Duration(24 * time.Hour).String(), "private code expiration")
	privateRetry 		= flag.String("R", time.Duration(30 * time.Second).String(), "private code regeneration interval")
)

func handleCodeRequest(p *burgoking.Pool, t string, w http.ResponseWriter, r *http.Request) {
	codeC, cancelC := make(chan string), make(chan struct{})
	go p.GetCode(codeC, cancelC)

	select {
	case code := <-codeC:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(code))
		logrus.WithFields(logrus.Fields{"code": code, "ip": realip.FromRequest(r), "type": t}).Info("Code used by user.")
	case <-r.Context().Done():
		cancelC <- struct{}{}
	}
}

func handlePublic(w http.ResponseWriter, r *http.Request) {
	handleCodeRequest(publicPool, "public", w, r)
}

func handlePrivate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *privateDirectKey {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	handleCodeRequest(privatePool, "private", w, r)
}

func handleDirect(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *privateDirectKey {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	code, err := burgoking.GenerateCode(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(code))
	logrus.WithFields(logrus.Fields{"code": code, "ip": r.RemoteAddr, "type": "direct"}).Info("Code used by user.")
}

func handleContact(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(*contact))
}

func parseDuration(s string) time.Duration {
	var duration time.Duration
	durationSec, err := strconv.Atoi(s)
	if err == nil {
		return time.Duration(durationSec) * time.Second
	} else {
		duration, err = time.ParseDuration(s)
		if err == nil {
			return duration
		}
	}

	return -1
}

func main() {
	flag.Parse()

	// Static files.
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Public pool.
	pool, err := burgoking.NewPool(*publicSize, parseDuration(*publicExpiration), parseDuration(*publicRetry))
	if err != nil {
		logrus.Fatal(err)
		return
	}
	publicPool = pool

	http.HandleFunc("/code", handlePublic)

	// Private pool.
	if *privateDirectKey != "" {
		pool, err = burgoking.NewPool(*privateSize, parseDuration(*privateExpiration), parseDuration(*privateRetry))
		if err != nil {
			logrus.Fatal(err)
			return
		}
		privatePool = pool

		http.HandleFunc("/private", 	handlePrivate)
		http.HandleFunc("/direct", 	handleDirect)

		logrus.Info("Private and direct code generation activated.")
	}

	// Contact address.
	if *contact != "" {
		http.HandleFunc("/contact", handleContact)
	}

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), nil))
}
