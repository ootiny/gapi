package gapi

import (
	"crypto/tls"
	"net/http"
)

var gMux = http.NewServeMux()

func ListenAndServe(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: gMux,
	}
	return server.ListenAndServe()
}

func ListenAndServeTLS(addr, certFile, keyFile string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: gMux,
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

func ListenAndServeTLSWithCert(addr string, certBytes []byte, keyBytes []byte) error {
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      addr,
		Handler:   gMux,
		TLSConfig: tlsConfig,
	}

	// Call ListenAndServeTLS on our custom server.
	// We pass empty strings for the certFile and keyFile because
	// the configuration is already provided in server.TLSConfig.
	return server.ListenAndServeTLS("", "")
}
