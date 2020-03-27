package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	"github.com/regalias/atlas-api/cache"
	"github.com/regalias/atlas-api/database"
	"github.com/regalias/atlas-api/logging"
)

type server struct {
	router           *httprouter.Router
	logger           *zerolog.Logger
	http             *http.Server
	databaseProvider database.Provider
	cacheProvider    cache.Provider
	indexSite        string
}

// Run runs the server instance
func Run(args []string) int {

	// Create logger
	lgr, err := logging.New("debug", "atlas-api", true)
	if err != nil {
		fmt.Printf("Oh noes! Something went horribly wrong!")
		panic(err)
	}

	r := httprouter.New()
	d, err := database.NewDDB(lgr, "atlas-table-main")
	if err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("Could not initialize database provider")
	}
	// TODO: grab table name from config
	if err := d.InitDatabase(); err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("Database or table was not found and could not create required resources")
	}

	// Create a redis cache provider
	c, err := cache.NewRedisProvider("127.0.0.1", 6379, nil)
	if err != nil {
		lgr.Fatal().Msg(err.Error())
	}

	// Create server context struct
	// TODO: pull from config
	s := server{
		router: r,
		logger: lgr,
		http: &http.Server{
			ReadHeaderTimeout: 20 * time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      2 * time.Minute,
			Addr:              ":8080",
			Handler:           r,
		},
		databaseProvider: d,
		cacheProvider:    c,
		indexSite:        "https://example.org",
	}

	s.routes(lgr)

	lgr.Info().Msg("Atlas server starting...")
	if err := s.http.ListenAndServe(); err != nil {
		lgr.Fatal().Err(err).Msg("API Startup failed")
	}

	return 0
}
