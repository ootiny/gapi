package __gapi_system__

import (
	"crypto/tls"
	"net/http"
)

var gMux = http.NewServeMux()

type IResponse interface {
}

type IRequest interface {
}

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

	return server.ListenAndServeTLS("", "")
}

func RegisterHandler(path string, handler func(IResponse, IRequest)) {
	gMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
}
