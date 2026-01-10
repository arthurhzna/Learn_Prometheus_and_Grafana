package main

import (
	"github.com/rs/zerolog/log"
)

func main () {
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // manage log level 
	log.Info().Msg("Hello, World!")
	log.Debug().Msg("Debug message")
	log.Warn().Msg("Warning message")
	log.Error().Msg("Error message")
	log.Fatal().Msg("Fatal message")
	// log.Panic().Msg("Panic message")
}