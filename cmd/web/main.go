package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

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

	publicSize       = flag.Int("n", 3, "public code pool size")
	publicExpiration = flag.Duration("d", 24*time.Hour, "public code expiration")
	publicRetry      = flag.Duration("r", 30*time.Second, "public code regeneration interval")

	privateDirectKey  = flag.String("k", "", "authorization token for private and direct code (disable if empty)")
	privateSize       = flag.Int("N", 1, "private code pool size")
	privateExpiration = flag.Duration("D", 24*time.Hour, "private code expiration")
	privateRetry      = flag.Duration("R", 30*time.Second, "private code regeneration interval")
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

func main() {
	flag.Parse()

	// Static files.
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Public pool.
	pool, err := burgoking.NewPool(*publicSize, *publicExpiration, *publicRetry)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	publicPool = pool

	http.HandleFunc("/code", handlePublic)

	// Private pool.
	if *privateDirectKey != "" {
		pool, err = burgoking.NewPool(*privateSize, *privateExpiration, *privateRetry)
		if err != nil {
			logrus.Fatal(err)
			return
		}
		privatePool = pool

		http.HandleFunc("/private", handlePrivate)
		http.HandleFunc("/direct", handleDirect)

		logrus.Info("Private and direct code generation activated.")
	}

	// Contact address.
	if *contact != "" {
		http.HandleFunc("/contact", handleContact)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
