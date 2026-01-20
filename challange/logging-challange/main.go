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

// package main

// import (
// 	"io"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// )

// func main() {

// 	f, err := os.OpenFile(
// 		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
// 	)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("unable to open log file")
// 	}
// 	defer f.Close()

// 	log.Logger = zerolog.New(f).With().Timestamp().Logger()

// 	gin.DefaultWriter = io.Discard
// 	gin.DefaultErrorWriter = io.Discard

// 	r := gin.New()
// 	r.Use(gin.Recovery())
// 	r.Use(RequestIDMiddleware())

// 	r.GET("/hello", func(c *gin.Context) {

// 		val, exists := c.Get("logger")
// 		var reqLogger zerolog.Logger
// 		if exists {
// 			reqLogger = val.(zerolog.Logger)
// 		} else {

// 			reqLogger = log.With().Logger()
// 		}

// 		reqLogger.Info().Msg("handler started")
// 		// contoh kerja
// 		time.Sleep(100 * time.Millisecond)
// 		reqLogger.Info().Msg("handler finished")

// 		c.JSON(200, gin.H{
// 			"message":    "hello",
// 			"request_id": c.GetString("request_id"),
// 		})
// 	})

// 	// jalankan server
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal().Err(err).Msg("server failed")
// 	}
// }

// func RequestIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := uuid.NewString()
// 		reqLogger := log.With().Str("request_id", id).Logger()

// 		c.Set("request_id", id)
// 		c.Set("logger", reqLogger)

// 		start := time.Now()
// 		reqLogger.Info().Msg("request started")

// 		c.Next() // jalankan handler

// 		latency := time.Since(start)
// 		status := c.Writer.Status()

// 		reqLogger.Info().
// 			Int("status", status).
// 			Dur("latency", latency).
// 			Msg("request completed")
// 	}
// // }

// package main

// import (
// 	"io"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promauto"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// )

// // Define Prometheus metrics
// var (
// 	// Counter for total HTTP requests
// 	httpRequestsTotal = promauto.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests",
// 		},
// 		[]string{"method", "endpoint", "status"},
// 	)

// 	// Histogram for request duration
// 	httpRequestDuration = promauto.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "HTTP request duration in seconds",
// 			Buckets: prometheus.DefBuckets, // Default buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
// 		},
// 		[]string{"method", "endpoint", "status"},
// 	)

// 	// Gauge for active requests
// 	httpRequestsInFlight = promauto.NewGauge(
// 		prometheus.GaugeOpts{
// 			Name: "http_requests_in_flight",
// 			Help: "Current number of HTTP requests being processed",
// 		},
// 	)
// )

// func main() {

// 	f, err := os.OpenFile(
// 		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
// 	)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("unable to open log file")
// 	}
// 	defer f.Close()

// 	log.Logger = zerolog.New(f).With().Timestamp().Logger()

// 	gin.DefaultWriter = io.Discard
// 	gin.DefaultErrorWriter = io.Discard

// 	r := gin.New()
// 	r.Use(gin.Recovery())
// 	r.Use(RequestIDMiddleware())
// 	r.Use(PrometheusMiddleware()) // Add Prometheus middleware

// 	r.GET("/hello", func(c *gin.Context) {

// 		val, exists := c.Get("logger")
// 		var reqLogger zerolog.Logger
// 		if exists {
// 			reqLogger = val.(zerolog.Logger)
// 		} else {
// 			reqLogger = log.With().Logger()
// 		}

// 		reqLogger.Info().Msg("handler started")
// 		// contoh kerja
// 		time.Sleep(100 * time.Millisecond)
// 		reqLogger.Info().Msg("handler finished")

// 		c.JSON(200, gin.H{
// 			"message":    "hello",
// 			"request_id": c.GetString("request_id"),
// 		})
// 	})

// 	// Prometheus metrics endpoint
// 	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

// 	// jalankan server
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal().Err(err).Msg("server failed")
// 	}
// }

// func RequestIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := uuid.NewString()
// 		reqLogger := log.With().Str("request_id", id).Logger()

// 		c.Set("request_id", id)
// 		c.Set("logger", reqLogger)

// 		start := time.Now()
// 		reqLogger.Info().Msg("request started")

// 		c.Next() // jalankan handler

// 		latency := time.Since(start)
// 		status := c.Writer.Status()

// 		reqLogger.Info().
// 			Int("status", status).
// 			Dur("latency", latency).
// 			Msg("request completed")
// 	}
// }

// // PrometheusMiddleware collects HTTP metrics
// func PrometheusMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Skip metrics endpoint to avoid recursion
// 		if c.Request.URL.Path == "/metrics" {
// 			c.Next()
// 			return
// 		}

// 		start := time.Now()

// 		// Increment in-flight requests
// 		httpRequestsInFlight.Inc()
// 		defer httpRequestsInFlight.Dec()

// 		c.Next()

// 		// Calculate duration
// 		duration := time.Since(start).Seconds()
// 		status := c.Writer.Status()
// 		method := c.Request.Method
// 		endpoint := c.FullPath()

// 		// If endpoint is empty (404), use the raw path
// 		if endpoint == "" {
// 			endpoint = c.Request.URL.Path
// 		}

// 		// Record metrics
// 		httpRequestsTotal.WithLabelValues(method, endpoint, string(rune(status))).Inc()
// 		httpRequestDuration.WithLabelValues(method, endpoint, string(rune(status))).Observe(duration)
// 	}
// }

package main

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Custom registry untuk kontrol penuh atas metrics
var (
	reg = prometheus.NewRegistry()

	// Counter untuk request count per endpoint (tanpa status)
	reqCountProcessed = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "The total number of processed by handler",
		},
		[]string{"method", "endpoint"},
	)

	// Counter untuk total HTTP requests (dengan status)
	httpRequestsTotal = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Histogram untuk request duration
	httpRequestDuration = promauto.With(reg).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	// Gauge untuk active requests
	httpRequestsInFlight = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

func main() {
	// Register Go runtime metrics collector
	reg.Register(collectors.NewGoCollector())

	f, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}
	defer f.Close()

	log.Logger = zerolog.New(f).With().Timestamp().Logger()

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())
	r.Use(PrometheusMiddleware())

	r.GET("/hello", func(c *gin.Context) {
		val, exists := c.Get("logger")
		var reqLogger zerolog.Logger
		if exists {
			reqLogger = val.(zerolog.Logger)
		} else {
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

	// Prometheus metrics endpoint dengan custom registry
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	// jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

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

// PrometheusMiddleware collects HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint to avoid recursion
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		endpoint := c.FullPath()

		// If endpoint is empty (404), use the raw path
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		// Convert status code to string properly
		statusStr := strconv.Itoa(status)

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, endpoint, statusStr).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint, statusStr).Observe(duration)

		// Custom counter tanpa status label
		reqCountProcessed.WithLabelValues(method, endpoint).Inc()
	}
}
