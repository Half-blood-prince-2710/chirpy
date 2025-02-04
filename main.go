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
type env struct {
	mode string
	jwtSecret string
}
type apiConfig struct {
	fileserverHits atomic.Int32
	db              *database.Queries // Initialize db as a pointer to database.Queries
	envi env
}

func main() {
	// Loading environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	mode:= os.Getenv("PLATFORM")
	slog.Info("Response: ","db_url", dbURL)

	// Database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Error connecting to database: ", "err",err)
	}
	defer db.Close()

	slog.Info("Database Connected")

	// Initialize the apiConfig object and assign the Queries object to db
	apiCfg := apiConfig{
		db: database.New(db),
		envi: env{
			mode: mode,
			jwtSecret: jwtSecret,
		} ,
	}

	// HTTP router setup
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.HandleFunc("GET /api/healthz", apiCfg.healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)


	//auth routes
	mux.HandleFunc("POST /api/login",apiCfg.loginHandler)
	mux.HandleFunc("POST /api/refresh",)

	// user routes
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	//chirps routes
	mux.HandleFunc("POST /api/chirps",apiCfg.authenticateMiddleware(apiCfg.createChirpHandler))
	mux.HandleFunc("GET /api/chirps",apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}",apiCfg.getChirpHandler)

	// Start HTTP server
	srv := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	slog.Info("server started...")
	log.Fatal(srv.ListenAndServe())
}





