package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}
	 mux := http.NewServeMux()
	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir("./")))))
	mux.HandleFunc("/healthz",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics",func(w http.ResponseWriter, r *http.Request) {
		fmt.Sprintf("Hits: ",apiCfg.fileserverHits)
	})
	 srv:= &http.Server{
		Handler: mux,
		Addr: ":8080",

	 }
	// log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
log.Fatal(srv.ListenAndServe())
	
}