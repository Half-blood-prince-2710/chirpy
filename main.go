package main

import (
	"log"
	"net/http"
	"sync/atomic"
)



func main() {
	 mux := http.NewServeMux()
	mux.Handle("/app/",http.StripPrefix("/app",http.FileServer(http.Dir("./"))))
	mux.HandleFunc("/healthz",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	 srv:= &http.Server{
		Handler: mux,
		Addr: ":8080",

	 }
	// log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
log.Fatal(srv.ListenAndServe())
	
}