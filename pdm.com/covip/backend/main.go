package main

import (
	//"fmt"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"pdm.com/covip/backend/helpers"
	"pdm.com/covip/backend/routes"
)

func main() {
	// Read Config file.
	helpers.ReadConfig()

	// Set up application arguments.
	addr := flag.String("addr", ":8080", "HTTPS network address")
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	r := routes.Routes()
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on port %s..", *addr)
	err := srv.ListenAndServeTLS(*certFile, *keyFile)
	log.Fatal(err)
}
