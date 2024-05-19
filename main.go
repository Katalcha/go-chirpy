package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Katalcha/go-chirpy/internal/database"
)

// SERVER CONFIG
const (
	LOCALHOST          string = "localhost"
	PORT               string = "8080"
	FILE_ROOT_PATH     string = "."
	FILE_DATABASE_PATH string = "database.json"
)

// ENDPOINTS
const (
	FILE_SERVER_PATH string = "/app/*"

	API_HEALTHZ string = "/api/healthz"

	API_CHIRPS         string = "/api/chirps"
	API_CHIRPS_ID      string = "/api/chirps/{chirpID}"
	API_VALIDATE_CHIRP string = "/api/validate_chirp"

	API_USERS    string = "/api/users"
	API_USERS_ID string = "/api/users/{userID}"

	API_LOGIN string = "/api/login"

	ADMIN_METRICS string = "/admin/metrics"
	API_RESET     string = "/api/reset"
)

// HTTP METHODS
const (
	GET  string = "GET "
	POST string = "POST "
)

// intern config struct to hold state
// fileServerHits - tracks the visitor count
type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	// reads or creates a ne DB ob server start, by checking for JSON-DB
	db, err := database.NewDB(FILE_DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// flag parsing for --debug, to delete database.json programatically
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	// create apiConfig for serverMetrics and in-memory DB
	apiCfg := apiConfig{fileServerHits: 0, DB: db}

	// create http server multiplexer
	serveMux := http.NewServeMux()

	// define file server
	fileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH))))
	serveMux.Handle(FILE_SERVER_PATH, fileServerHandler)

	// let multiplexer handle specific endpoints
	serveMux.HandleFunc(GET+API_HEALTHZ, healthzHandler) // get readiness on GET /api/healthz

	serveMux.HandleFunc(GET+API_CHIRPS, apiCfg.getChirpsHandler)       // gets all chirps in database on GET /api/chirps
	serveMux.HandleFunc(POST+API_CHIRPS, apiCfg.createChirpHandler)    // posts a new chirp with inbund validation on POST /api/chirps
	serveMux.HandleFunc(GET+API_CHIRPS_ID, apiCfg.getChirpByIdHandler) // gets a specific chirp in database by id on GET /api/chirps/{chirpID}

	serveMux.HandleFunc(GET+API_USERS, apiCfg.getUsersHandler)       // gets all users in database on GET /api/users
	serveMux.HandleFunc(GET+API_USERS_ID, apiCfg.getUserByIdHandler) // gets a specific user in database by id on GET /api/users/{userID}
	serveMux.HandleFunc(POST+API_LOGIN, apiCfg.loginUserHandler)
	serveMux.HandleFunc(POST+API_USERS, apiCfg.createUserHandler) // creates a new user on POST /api/users

	serveMux.HandleFunc(GET+ADMIN_METRICS, apiCfg.metricsHandler)  // get visitor count metrics on GET /admin/metrics
	serveMux.HandleFunc(GET+API_RESET, apiCfg.metricsResetHandler) // resets visitor cound metrics on GET /api/reset
	// serveMux.HandleFunc(POST+API_VALIDATE_CHIRP, validateChirpHandler) // old: validiates a posted chirp on structure and profanity on POST /api/validate_chirp

	// create http.Server object with configured serveMux
	httpServer := &http.Server{Addr: LOCALHOST + ":" + PORT, Handler: serveMux}

	// log info, start server, inform on fatal or close
	log.Printf("Serving Yo Mama from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	log.Fatal(httpServer.ListenAndServe())
}
