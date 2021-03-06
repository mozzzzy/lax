package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mkusaka/lax/client"
)

func main() {
	// TODO: serve via H2
	// TODO: allow auto restart like HUP USR2, or http request like /restart
	http.Handle("/", http.HandlerFunc(proxyHander))
	fmt.Println("served at http://localhost:300")
	http.ListenAndServe(":300", nil)
}

func proxyHander(w http.ResponseWriter, r *http.Request) {
	c := client.NewClient(10 * time.Second)
	// TODO: not ruled pattern makes panic
	res, err := c.ProxyRequest(r)

	if err != nil {
		// TODO: log store from worker
		// TODO: error handling. add retry function
		// TODO: send error report to primary server
		http.Error(w, http.StatusText(http.StatusBadRequest)+" :"+err.Error(), http.StatusBadRequest)
		return
	}

	// res.Header returns map[string][]string
	// TODO: following code valid? (write all header from origin server to response from proxy server)
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("invalid encoding: %s", err)
		}
		w.Write(body)
	}
}
