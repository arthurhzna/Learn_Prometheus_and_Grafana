// package main

// import (
// 	"os"

// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

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

// func main() {
// 	f, err := os.OpenFile(
// 		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
// 	)

// 	if err != nil {
// 		log.Fatal().Err(err).Msg("unable to open log file")
// 	}
// 	defer f.Close()

// 	log.Logger = zerolog.New(f).With().Timestamp().Logger()
// 	log.Debug().Msg("Debug message")

// }

// HTTP Request Logging

package main

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// buka file log (append/create)
	f, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}
	defer f.Close()

	// set global logger ke file (JSON structured)
	log.Logger = zerolog.New(f).With().Timestamp().Logger()

	// hindari Gin menulis ke stdout (agar Fluent Bit baca hanya file)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())

	// contoh handler
	r.GET("/hello", func(c *gin.Context) {
		// ambil logger yang diset oleh middleware
		val, exists := c.Get("logger")
		var reqLogger zerolog.Logger
		if exists {
			reqLogger = val.(zerolog.Logger)
		} else {
			// fallback ke global logger dengan request_id kosong
			reqLogger = log.With().Logger()
		}

		reqLogger.Info().Msg("handler started")
		// contoh kerja
		time.Sleep(100 * time.Millisecond)
		reqLogger.Info().Msg("handler finished")

		c.JSON(200, gin.H{
			"message":    "hello",
			"request_id": c.GetString("request_id"),
		})
	})

	// jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

// RequestIDMiddleware men-set UUID per request dan menulis log sebelum & sesudah handler
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		reqLogger := log.With().Str("request_id", id).Logger()

		c.Set("request_id", id)
		c.Set("logger", reqLogger)

		start := time.Now()
		reqLogger.Info().Msg("request started")

		c.Next() // jalankan handler

		latency := time.Since(start)
		status := c.Writer.Status()

		reqLogger.Info().
			Int("status", status).
			Dur("latency", latency).
			Msg("request completed")
	}
}
