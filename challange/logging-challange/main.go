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

// package main

// import (
// 	"io"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/collectors"
// 	"github.com/prometheus/client_golang/prometheus/promauto"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// )

// // Custom registry untuk kontrol penuh atas metrics
// var (
// 	reg = prometheus.NewRegistry()

// 	// Counter untuk request count per endpoint (tanpa status)
// 	reqCountProcessed = promauto.With(reg).NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_request_count",
// 			Help: "The total number of processed by handler",
// 		},
// 		[]string{"method", "endpoint"},
// 	)

// 	// Counter untuk total HTTP requests (dengan status)
// 	httpRequestsTotal = promauto.With(reg).NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests",
// 		},
// 		[]string{"method", "endpoint", "status"},
// 	)

// 	// Histogram untuk request duration
// 	httpRequestDuration = promauto.With(reg).NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "HTTP request duration in seconds",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"method", "endpoint", "status"},
// 	)

// 	// Gauge untuk active requests
// 	httpRequestsInFlight = promauto.With(reg).NewGauge(
// 		prometheus.GaugeOpts{
// 			Name: "http_requests_in_flight",
// 			Help: "Current number of HTTP requests being processed",
// 		},
// 	)
// )

// func main() {
// 	// Register Go runtime metrics collector
// 	reg.Register(collectors.NewGoCollector())

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
// 	r.Use(PrometheusMiddleware())

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

// 	// Prometheus metrics endpoint dengan custom registry
// 	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

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

// 		// Convert status code to string properly
// 		statusStr := strconv.Itoa(status)

// 		// Record metrics
// 		httpRequestsTotal.WithLabelValues(method, endpoint, statusStr).Inc()
// 		httpRequestDuration.WithLabelValues(method, endpoint, statusStr).Observe(duration)

// 		// Custom counter tanpa status label
// 		reqCountProcessed.WithLabelValues(method, endpoint).Inc()
// 	}
// }

// package main

// import (
// 	"context"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
// 	"go.opentelemetry.io/otel/trace"
// )

// var tracer trace.Tracer

// func main() {

// 	shutdown, err := initTracer()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to initialize tracer")
// 	}
// 	defer func() {
// 		if err := shutdown(context.Background()); err != nil {
// 			log.Fatal().Err(err).Msg("failed to shutdown tracer")
// 		}
// 	}()

// 	f, err := os.OpenFile(
// 		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
// 	)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("unable to open log file")
// 	}
// 	defer f.Close()

// 	log.Logger = zerolog.New(f).With().Timestamp().Logger()

// 	// gin.DefaultWriter = io.Discard
// 	// gin.DefaultErrorWriter = io.Discard

// 	r := gin.New()
// 	r.Use(gin.Recovery())
// 	r.Use(RequestIDMiddleware())
// 	r.Use(OtelMiddleware())

// 	r.GET("/hello", helloHandler)
// 	r.GET("/slow", slowHandler)
// 	r.GET("/error", errorHandler)

// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal().Err(err).Msg("server failed")
// 	}
// }

// func initTracer() (func(context.Context) error, error) {
// 	exporter, err := otlptracehttp.New(context.Background(),
// 		otlptracehttp.WithEndpoint("localhost:4318"),
// 		otlptracehttp.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := resource.New(context.Background(),
// 		resource.WithAttributes(
// 			semconv.ServiceNameKey.String("logging-challenge-service"),
// 			semconv.ServiceVersionKey.String("1.0.0"),
// 			attribute.String("environment", "development"),
// 		),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exporter),
// 		sdktrace.WithResource(res),
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 	)

// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	))

// 	tracer = tp.Tracer("gin-server")

// 	return tp.Shutdown, nil
// }

// func OtelMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()

// 		spanName := c.Request.Method + " " + c.FullPath()
// 		if c.FullPath() == "" {
// 			spanName = c.Request.Method + " " + c.Request.URL.Path
// 		}

// 		ctx, span := tracer.Start(ctx, spanName,
// 			trace.WithSpanKind(trace.SpanKindServer),
// 		)
// 		defer span.End()

// 		span.SetAttributes(
// 			semconv.HTTPMethodKey.String(c.Request.Method),
// 			semconv.HTTPRouteKey.String(c.FullPath()),
// 			semconv.HTTPTargetKey.String(c.Request.URL.Path),
// 			semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
// 			attribute.String("http.client_ip", c.ClientIP()),
// 			attribute.String("request_id", c.GetString("request_id")),
// 		)

// 		c.Request = c.Request.WithContext(ctx)

// 		start := time.Now()
// 		c.Next()
// 		duration := time.Since(start)

// 		status := c.Writer.Status()
// 		span.SetAttributes(
// 			semconv.HTTPStatusCodeKey.Int(status),
// 			attribute.Int64("http.duration_ms", duration.Milliseconds()),
// 		)

// 		if len(c.Errors) > 0 {
// 			span.RecordError(c.Errors.Last())
// 		}
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

// 		c.Next()

// 		latency := time.Since(start)
// 		status := c.Writer.Status()

// 		reqLogger.Info().
// 			Int("status", status).
// 			Dur("latency", latency).
// 			Msg("request completed")
// 	}
// }

// func helloHandler(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	val, exists := c.Get("logger")
// 	var reqLogger zerolog.Logger
// 	if exists {
// 		reqLogger = val.(zerolog.Logger)
// 	} else {
// 		reqLogger = log.With().Logger()
// 	}

// 	_, span := tracer.Start(ctx, "helloHandler.processing")
// 	defer span.End()

// 	reqLogger.Info().Msg("handler started")

// 	time.Sleep(100 * time.Millisecond)

// 	span.AddEvent("work completed")
// 	reqLogger.Info().Msg("handler finished")

// 	c.JSON(200, gin.H{
// 		"message":    "hello",
// 		"request_id": c.GetString("request_id"),
// 	})
// }

// func slowHandler(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	val, exists := c.Get("logger")
// 	var reqLogger zerolog.Logger
// 	if exists {
// 		reqLogger = val.(zerolog.Logger)
// 	} else {
// 		reqLogger = log.With().Logger()
// 	}

// 	_, span1 := tracer.Start(ctx, "slowHandler.database_query")
// 	reqLogger.Info().Msg("querying database")
// 	time.Sleep(200 * time.Millisecond)
// 	span1.SetAttributes(attribute.String("db.query", "SELECT * FROM users"))
// 	span1.End()

// 	_, span2 := tracer.Start(ctx, "slowHandler.external_api_call")
// 	reqLogger.Info().Msg("calling external API")
// 	time.Sleep(300 * time.Millisecond)
// 	span2.SetAttributes(attribute.String("http.url", "https://api.example.com"))
// 	span2.End()

// 	reqLogger.Info().Msg("slow handler finished")

// 	c.JSON(200, gin.H{
// 		"message":  "slow operation completed",
// 		"duration": "500ms",
// 	})
// }

// func errorHandler(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	val, exists := c.Get("logger")
// 	var reqLogger zerolog.Logger
// 	if exists {
// 		reqLogger = val.(zerolog.Logger)
// 	} else {
// 		reqLogger = log.With().Logger()
// 	}

// 	_, span := tracer.Start(ctx, "errorHandler.processing")
// 	defer span.End()

// 	reqLogger.Error().Msg("simulated error occurred")

// 	span.RecordError(gin.Error{
// 		Err:  gin.Error{}.Err,
// 		Type: gin.ErrorTypePrivate,
// 		Meta: "simulated error",
// 	})
// 	span.SetAttributes(attribute.Bool("error", true))

// 	c.JSON(500, gin.H{
// 		"error": "something went wrong",
// 	})
// }

// -------------------------------------------- prometheus exporter and otel tracer --------------------------------------------
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	meter  metric.Meter

	httpRequestCounter  metric.Int64Counter
	httpRequestDuration metric.Float64Histogram
	httpRequestsActive  metric.Int64UpDownCounter
)

func main() {
	// Initialize OpenTelemetry (traces + metrics)
	shutdownTrace, metricsHandler, err := initObservability()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize observability")
	}
	defer func() {
		if err := shutdownTrace(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown tracer")
		}
	}()

	f, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}
	defer f.Close()

	log.Logger = zerolog.New(f).With().Timestamp().Logger()

	// Uncomment to see Gin output
	// gin.DefaultWriter = os.Stdout
	// gin.DefaultErrorWriter = os.Stderr

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())
	r.Use(OtelMiddleware())

	r.GET("/hello", helloHandler)
	r.GET("/slow", slowHandler)
	r.GET("/error", errorHandler)
	r.GET("/external", externalAPIHandler) // New: Demo otelhttp for external calls

	// Prometheus metrics endpoint - PERBAIKAN DI SINI
	r.GET("/metrics", func(c *gin.Context) {
		metricsHandler.ServeHTTP(c.Writer, c.Request)
	})

	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("📊 Traces: http://localhost:16686 (Jaeger UI)")
	fmt.Println("📈 Metrics: http://localhost:8080/metrics")
	fmt.Println("📝 Logs: logs/app.log")
	fmt.Println("\nEndpoints:")
	fmt.Println("  - GET /hello")
	fmt.Println("  - GET /slow")
	fmt.Println("  - GET /error")
	fmt.Println("  - GET /external (demo otelhttp)")
	fmt.Println("  - GET /metrics (Prometheus)")

	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

// PERBAIKAN: Return http.Handler instead of *prometheus.Exporter
func initObservability() (func(context.Context) error, http.Handler, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("logging-challenge-service"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	// === TRACES (push to OTLP/Jaeger) ===
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("gin-server")

	// === METRICS (Prometheus exporter) ===
	promExporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)
	meter = mp.Meter("gin-server")

	// Create metrics instruments
	httpRequestCounter, err = meter.Int64Counter(
		"http_server_request_count",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, nil, err
	}

	httpRequestDuration, err = meter.Float64Histogram(
		"http_server_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, nil, err
	}

	httpRequestsActive, err = meter.Int64UpDownCounter(
		"http_server_requests_active",
		metric.WithDescription("Number of active HTTP requests"),
	)
	if err != nil {
		return nil, nil, err
	}

	// Set propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// PERBAIKAN: Use promhttp.Handler() to expose metrics
	// The OpenTelemetry Prometheus exporter registers metrics with the default Prometheus registry
	// So we can use the standard Prometheus HTTP handler
	return tp.Shutdown, promhttp.Handler(), nil
}

func OtelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPRouteKey.String(c.FullPath()),
			semconv.HTTPTargetKey.String(c.Request.URL.Path),
			semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
			attribute.String("http.client_ip", c.ClientIP()),
			attribute.String("request_id", c.GetString("request_id")),
		)

		c.Request = c.Request.WithContext(ctx)

		// Increment active requests
		httpRequestsActive.Add(ctx, 1)
		defer httpRequestsActive.Add(ctx, -1)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		status := c.Writer.Status()
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(status),
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
		)

		// Record metrics
		attrs := metric.WithAttributes(
			attribute.String("method", c.Request.Method),
			attribute.String("route", c.FullPath()),
			attribute.Int("status", status),
		)

		httpRequestCounter.Add(ctx, 1, attrs)
		httpRequestDuration.Record(ctx, duration.Seconds(), attrs)

		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last())
		}
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

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		reqLogger.Info().
			Int("status", status).
			Dur("latency", latency).
			Msg("request completed")
	}
}

func helloHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	_, span := tracer.Start(ctx, "helloHandler.processing")
	defer span.End()

	reqLogger.Info().Msg("handler started")
	time.Sleep(100 * time.Millisecond)
	span.AddEvent("work completed")
	reqLogger.Info().Msg("handler finished")

	c.JSON(200, gin.H{
		"message":    "hello",
		"request_id": c.GetString("request_id"),
	})
}

func slowHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	_, span1 := tracer.Start(ctx, "slowHandler.database_query")
	reqLogger.Info().Msg("querying database")
	time.Sleep(200 * time.Millisecond)
	span1.SetAttributes(attribute.String("db.query", "SELECT * FROM users"))
	span1.End()

	_, span2 := tracer.Start(ctx, "slowHandler.external_api_call")
	reqLogger.Info().Msg("calling external API")
	time.Sleep(300 * time.Millisecond)
	span2.SetAttributes(attribute.String("http.url", "https://api.example.com"))
	span2.End()

	reqLogger.Info().Msg("slow handler finished")

	c.JSON(200, gin.H{
		"message":  "slow operation completed",
		"duration": "500ms",
	})
}

func errorHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	_, span := tracer.Start(ctx, "errorHandler.processing")
	defer span.End()

	reqLogger.Error().Msg("simulated error occurred")

	span.RecordError(gin.Error{
		Err:  gin.Error{}.Err,
		Type: gin.ErrorTypePrivate,
		Meta: "simulated error",
	})
	span.SetAttributes(attribute.Bool("error", true))

	c.JSON(500, gin.H{
		"error": "something went wrong",
	})
}

// externalAPIHandler demonstrates using otelhttp for automatic HTTP client instrumentation
func externalAPIHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	_, span := tracer.Start(ctx, "externalAPIHandler.processing")
	defer span.End()

	reqLogger.Info().Msg("making external API calls with otelhttp")

	// Create HTTP client with otelhttp instrumentation
	// This will automatically create spans for HTTP calls
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Example 1: Call to JSONPlaceholder API
	reqLogger.Info().Msg("calling jsonplaceholder API")
	req1, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		reqLogger.Error().Err(err).Msg("failed to create request")
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	resp1, err := client.Do(req1)
	if err != nil {
		reqLogger.Error().Err(err).Msg("failed to call jsonplaceholder API")
		span.RecordError(err)
		c.JSON(500, gin.H{"error": "failed to call external API"})
		return
	}
	defer resp1.Body.Close()

	span.AddEvent("jsonplaceholder API called successfully")
	reqLogger.Info().Int("status", resp1.StatusCode).Msg("jsonplaceholder API response")

	// Example 2: Call to another endpoint (will create another child span)
	reqLogger.Info().Msg("calling another API")
	req2, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/users/1", nil)
	if err != nil {
		reqLogger.Error().Err(err).Msg("failed to create request")
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	resp2, err := client.Do(req2)
	if err != nil {
		reqLogger.Error().Err(err).Msg("failed to call users API")
		span.RecordError(err)
		c.JSON(500, gin.H{"error": "failed to call external API"})
		return
	}
	defer resp2.Body.Close()

	span.AddEvent("users API called successfully")
	reqLogger.Info().Int("status", resp2.StatusCode).Msg("users API response")

	// Add span attributes
	span.SetAttributes(
		attribute.Int("external.api1.status", resp1.StatusCode),
		attribute.Int("external.api2.status", resp2.StatusCode),
		attribute.Bool("success", true),
	)

	reqLogger.Info().Msg("external API handler finished")

	c.JSON(200, gin.H{
		"message": "external API calls completed",
		"api1": gin.H{
			"url":    "https://jsonplaceholder.typicode.com/posts/1",
			"status": resp1.StatusCode,
		},
		"api2": gin.H{
			"url":    "https://jsonplaceholder.typicode.com/users/1",
			"status": resp2.StatusCode,
		},
		"request_id": c.GetString("request_id"),
	})
}
