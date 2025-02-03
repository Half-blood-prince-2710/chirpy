package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/half-blood-prince-2710/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries  *database.Queries
}

func main() {
	// loading environment variables 
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	// database connection
	db, err := sql.Open("postgres", dbURL)
	if err!=nil {
		slog.Error("Error connecting database")
	}
	defer db.Close()

	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.HandleFunc("GET /api/healthz", apiCfg.healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp",apiCfg.validateChirpHandler)
	mux.HandleFunc("POST /api/users",apiC)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	// log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
