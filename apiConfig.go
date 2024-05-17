package main

import (
	"github.com/Katalcha/go-chirpy/internal/database"
)

// intern config struct to hold state
// fileServerHits - tracks the visitor count
type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}
