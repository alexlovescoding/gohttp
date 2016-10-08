package main

import (
	"os"
	"net/http"
	"log"
)

func redirect(w http.ResponseWriter, req *http.Request) {
	log.Println("Redirecting to https...")
	http.Redirect(w, req,
		"https://" + req.Host + req.URL.String(),
		http.StatusMovedPermanently)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/public-html")))

	errchan := make(chan error)

	if _, err := os.Stat("/tls/current.chain"); err != nil && os.IsNotExist(err) {
		log.Println("Serving http...")
		go func() {
			if err := http.ListenAndServe(":80", nil); err != nil {
				errchan <- err
			}
		}()
	} else {
		log.Println("Serving https...")
		go func() {
			if err := http.ListenAndServeTLS(":443", "/tls/current.chain", "/tls/current.key", nil); err != nil {
				errchan <- err
			}
		}()

		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirect)); err != nil {
				errchan <- err
			}
		}()
	}

	select {
	case err := <- errchan:
		log.Fatalln("Could not start serving service due to (error: %s)", err)
	}
}
