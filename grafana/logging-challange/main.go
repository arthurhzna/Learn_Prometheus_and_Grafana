package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// func main () {
// 	zerolog.SetGlobalLevel(zerolog.InfoLevel) // manage log level
// 	log.Info().Msg("Hello, World!")
// 	log.Debug().Msg("Debug message")
// 	log.Warn().Msg("Warning message")
// 	log.Error().Msg("Error message")
// 	log.Fatal().Msg("Fatal message")
// 	// log.Panic().Msg("Panic message")
// }

// structured logging
// func main() {

// 	log := log.With().
// 		Str("booking_id", "1234567890").
// 		Float64("price", 100.00).
// 		Str("currency", "USD").
// 		Logger()

// 	log.Debug().
// 		Int("user_age", 20).
// 		Msg("Booking created")

// 	err := fmt.Errorf("Ticket not found")
// 	log.Error().Err(err).Msgf("Booking is failed %s", "man")

// }

//log-correlation-id
// track everything using uuid to make it easier to trace

// func makeBooking(ctx context.Context) {
// 	log.Ctx(ctx).Info().Msg("Creating a Booking")
// }

// func main() {

// 	log := log.With().
// 		Str("booking_id", "1234567890").
// 		Float64("price", 100.00).
// 		Str("currency", "USD").
// 		Logger()

// 	ctx := log.WithContext(context.Background())

// 	makeBooking(ctx)

// 	log.Debug().
// 		Int("user_age", 20).
// 		Msg("Booking created")
// }

// {"level":"info","booking_id":"1234567890","price":100,"currency":"USD","time":"2026-01-10T13:08:50+07:00","message":"Creating a Booking"}
// {"level":"debug","booking_id":"1234567890","price":100,"currency":"USD","user_age":20,"time":"2026-01-10T13:08:50+07:00","message":"Booking created"}
// same booking_id but different message level with using context

// in development we write logs to stdout/ console and file log
// in production we write logs to stdout/ console or file log, usually use stoud/console for easier collect by local agent or kubernetes agent

//wriet log into file log using zerolog

func main() {
	f, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}

	defer f.Close()

	log.Logger = zerolog.New(f).With().Timestamp().Logger()
	log.Debug().Msg("Debug message")

}
