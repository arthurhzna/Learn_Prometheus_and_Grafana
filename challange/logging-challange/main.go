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

// package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// 	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/exporters/prometheus"
// 	"go.opentelemetry.io/otel/metric"
// 	"go.opentelemetry.io/otel/propagation"
// 	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
// 	"go.opentelemetry.io/otel/trace"
// )

// var (
// 	tracer trace.Tracer
// 	meter  metric.Meter

// 	httpRequestCounter  metric.Int64Counter
// 	httpRequestDuration metric.Float64Histogram
// 	httpRequestsActive  metric.Int64UpDownCounter
// )

// func main() {
// 	// Initialize OpenTelemetry (traces + metrics)
// 	shutdownTrace, metricsHandler, err := initObservability()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to initialize observability")
// 	}
// 	defer func() {
// 		if err := shutdownTrace(context.Background()); err != nil {
// 			log.Error().Err(err).Msg("failed to shutdown tracer")
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

// 	// Uncomment to see Gin output
// 	// gin.DefaultWriter = os.Stdout
// 	// gin.DefaultErrorWriter = os.Stderr

// 	r := gin.New()
// 	r.Use(gin.Recovery())
// 	r.Use(RequestIDMiddleware())
// 	r.Use(OtelMiddleware())

// 	r.GET("/hello", helloHandler)
// 	r.GET("/slow", slowHandler)
// 	r.GET("/error", errorHandler)
// 	r.GET("/external", externalAPIHandler) // New: Demo otelhttp for external calls

// 	// Prometheus metrics endpoint - PERBAIKAN DI SINI
// 	r.GET("/metrics", func(c *gin.Context) {
// 		metricsHandler.ServeHTTP(c.Writer, c.Request)
// 	})

// 	fmt.Println("🚀 Server starting on http://localhost:8080")
// 	fmt.Println("📊 Traces: http://localhost:16686 (Jaeger UI)")
// 	fmt.Println("📈 Metrics: http://localhost:8080/metrics")
// 	fmt.Println("📝 Logs: logs/app.log")
// 	fmt.Println("\nEndpoints:")
// 	fmt.Println("  - GET /hello")
// 	fmt.Println("  - GET /slow")
// 	fmt.Println("  - GET /error")
// 	fmt.Println("  - GET /external (demo otelhttp)")
// 	fmt.Println("  - GET /metrics (Prometheus)")

// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal().Err(err).Msg("server failed")
// 	}
// }

// // PERBAIKAN: Return http.Handler instead of *prometheus.Exporter
// func initObservability() (func(context.Context) error, http.Handler, error) {
// 	ctx := context.Background()

// 	res, err := resource.New(ctx,
// 		resource.WithAttributes(
// 			semconv.ServiceNameKey.String("logging-challenge-service"),
// 			semconv.ServiceVersionKey.String("1.0.0"),
// 			attribute.String("environment", "development"),
// 		),
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// === TRACES (push to OTLP/Jaeger) ===
// 	traceExporter, err := otlptracehttp.New(ctx,
// 		otlptracehttp.WithEndpoint("localhost:4318"),
// 		otlptracehttp.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(traceExporter),
// 		sdktrace.WithResource(res),
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 	)

// 	otel.SetTracerProvider(tp)
// 	tracer = tp.Tracer("gin-server")

// 	// === METRICS (Prometheus exporter) ===
// 	promExporter, err := prometheus.New()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	mp := sdkmetric.NewMeterProvider(
// 		sdkmetric.WithReader(promExporter),
// 		sdkmetric.WithResource(res),
// 	)

// 	otel.SetMeterProvider(mp)
// 	meter = mp.Meter("gin-server")

// 	// Create metrics instruments
// 	httpRequestCounter, err = meter.Int64Counter(
// 		"http_server_request_count",
// 		metric.WithDescription("Total number of HTTP requests"),
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	httpRequestDuration, err = meter.Float64Histogram(
// 		"http_server_request_duration_seconds",
// 		metric.WithDescription("HTTP request duration in seconds"),
// 		metric.WithUnit("s"),
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	httpRequestsActive, err = meter.Int64UpDownCounter(
// 		"http_server_requests_active",
// 		metric.WithDescription("Number of active HTTP requests"),
// 	)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Set propagators
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	))

// 	// PERBAIKAN: Use promhttp.Handler() to expose metrics
// 	// The OpenTelemetry Prometheus exporter registers metrics with the default Prometheus registry
// 	// So we can use the standard Prometheus HTTP handler
// 	return tp.Shutdown, promhttp.Handler(), nil
// }

// func OtelMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Skip metrics endpoint
// 		if c.Request.URL.Path == "/metrics" {
// 			c.Next()
// 			return
// 		}

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

// 		// Increment active requests
// 		httpRequestsActive.Add(ctx, 1)
// 		defer httpRequestsActive.Add(ctx, -1)

// 		start := time.Now()
// 		c.Next()
// 		duration := time.Since(start)

// 		status := c.Writer.Status()
// 		span.SetAttributes(
// 			semconv.HTTPStatusCodeKey.Int(status),
// 			attribute.Int64("http.duration_ms", duration.Milliseconds()),
// 		)

// 		// Record metrics
// 		attrs := metric.WithAttributes(
// 			attribute.String("method", c.Request.Method),
// 			attribute.String("route", c.FullPath()),
// 			attribute.Int("status", status),
// 		)

// 		httpRequestCounter.Add(ctx, 1, attrs)
// 		httpRequestDuration.Record(ctx, duration.Seconds(), attrs)

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

// func externalAPIHandler(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	val, exists := c.Get("logger")
// 	var reqLogger zerolog.Logger
// 	if exists {
// 		reqLogger = val.(zerolog.Logger)
// 	} else {
// 		reqLogger = log.With().Logger()
// 	}

// 	_, span := tracer.Start(ctx, "externalAPIHandler.processing")
// 	defer span.End()

// 	reqLogger.Info().Msg("making external API calls with otelhttp")

// 	client := http.Client{
// 		Transport: otelhttp.NewTransport(http.DefaultTransport),
// 	}

// 	reqLogger.Info().Msg("calling jsonplaceholder API")
// 	req1, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
// 	if err != nil {
// 		reqLogger.Error().Err(err).Msg("failed to create request")
// 		c.JSON(500, gin.H{"error": "failed to create request"})
// 		return
// 	}

// 	resp1, err := client.Do(req1)
// 	if err != nil {
// 		reqLogger.Error().Err(err).Msg("failed to call jsonplaceholder API")
// 		span.RecordError(err)
// 		c.JSON(500, gin.H{"error": "failed to call external API"})
// 		return
// 	}
// 	defer resp1.Body.Close()

// 	span.AddEvent("jsonplaceholder API called successfully")
// 	reqLogger.Info().Int("status", resp1.StatusCode).Msg("jsonplaceholder API response")

// 	reqLogger.Info().Msg("calling another API")
// 	req2, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/users/1", nil)
// 	if err != nil {
// 		reqLogger.Error().Err(err).Msg("failed to create request")
// 		c.JSON(500, gin.H{"error": "failed to create request"})
// 		return
// 	}

// 	resp2, err := client.Do(req2)
// 	if err != nil {
// 		reqLogger.Error().Err(err).Msg("failed to call users API")
// 		span.RecordError(err)
// 		c.JSON(500, gin.H{"error": "failed to call external API"})
// 		return
// 	}
// 	defer resp2.Body.Close()

// 	span.AddEvent("users API called successfully")
// 	reqLogger.Info().Int("status", resp2.StatusCode).Msg("users API response")

// 	span.SetAttributes(
// 		attribute.Int("external.api1.status", resp1.StatusCode),
// 		attribute.Int("external.api2.status", resp2.StatusCode),
// 		attribute.Bool("success", true),
// 	)

// 	reqLogger.Info().Msg("external API handler finished")

// 	c.JSON(200, gin.H{
// 		"message": "external API calls completed",
// 		"api1": gin.H{
// 			"url":    "https://jsonplaceholder.typicode.com/posts/1",
// 			"status": resp1.StatusCode,
// 		},
// 		"api2": gin.H{
// 			"url":    "https://jsonplaceholder.typicode.com/users/1",
// 			"status": resp2.StatusCode,
// 		},
// 		"request_id": c.GetString("request_id"),
// 	})
// }

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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())
	r.Use(OtelMiddleware())

	r.GET("/hello", helloHandler)
	r.GET("/slow", slowHandler)
	r.GET("/error", errorHandler)
	r.GET("/external", externalAPIHandler)
	r.GET("/propagation", propagationDemoHandler) // Demo context propagation with otelhttp

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
	fmt.Println("  - GET /propagation (demo context propagation)")
	fmt.Println("  - GET /metrics (Prometheus)")

	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

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

	// === TRACES (push to OTLP/Jaeger via gRPC) ===
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
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

	// Add span attributes for better observability
	span.SetAttributes(
		attribute.String("handler.name", "hello"),
		attribute.String("handler.type", "simple"),
		attribute.String("request_id", c.GetString("request_id")),
	)

	// Add event: handler started
	span.AddEvent("handler.started", trace.WithAttributes(
		attribute.String("state", "processing"),
	))
	reqLogger.Info().Msg("handler started")

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	// Add event: work completed
	span.AddEvent("work.completed", trace.WithAttributes(
		attribute.Int64("duration_ms", 100),
		attribute.String("state", "success"),
	))
	reqLogger.Info().Msg("handler finished")

	// Mark span as successful
	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.String("response.status", "ok"),
	)

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

	// Database query span with detailed attributes
	ctx, span1 := tracer.Start(ctx, "slowHandler.database_query",
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "SELECT"),
		),
	)
	span1.AddEvent("query.started", trace.WithAttributes(
		attribute.String("state", "executing"),
	))
	reqLogger.Info().Msg("querying database")

	time.Sleep(200 * time.Millisecond)

	span1.SetAttributes(
		attribute.String("db.query", "SELECT * FROM users"),
		attribute.String("db.table", "users"),
		attribute.Int("db.rows_returned", 42),
		attribute.Int64("db.duration_ms", 200),
		attribute.Bool("db.cache_hit", false),
	)
	span1.AddEvent("query.completed", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("rows", 42),
	))
	span1.End()

	// External API call span with detailed attributes
	ctx, span2 := tracer.Start(ctx, "slowHandler.external_api_call",
		trace.WithAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.url", "https://api.example.com"),
		),
	)
	span2.AddEvent("api_call.started", trace.WithAttributes(
		attribute.String("state", "connecting"),
	))
	reqLogger.Info().Msg("calling external API")

	time.Sleep(300 * time.Millisecond)

	span2.SetAttributes(
		attribute.Int("http.status_code", 200),
		attribute.Int64("http.duration_ms", 300),
		attribute.Int("http.response_size_bytes", 1024),
		attribute.Bool("http.cached", false),
	)
	span2.AddEvent("api_call.completed", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("status_code", 200),
	))
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

	_, span := tracer.Start(ctx, "errorHandler.processing",
		trace.WithAttributes(
			attribute.String("handler.name", "error"),
			attribute.String("handler.type", "error_simulation"),
		),
	)
	defer span.End()

	// Add event: error processing started
	span.AddEvent("error.processing_started", trace.WithAttributes(
		attribute.String("state", "processing"),
	))

	// Simulate error scenario
	err := fmt.Errorf("simulated database connection error")
	reqLogger.Error().Err(err).Msg("simulated error occurred")

	// Record error with detailed attributes
	span.RecordError(err, trace.WithAttributes(
		attribute.String("error.type", "DatabaseConnectionError"),
		attribute.String("error.message", err.Error()),
		attribute.String("error.severity", "high"),
		attribute.String("error.component", "database"),
	))

	// Add error event
	span.AddEvent("error.occurred", trace.WithAttributes(
		attribute.String("state", "failed"),
		attribute.String("error.type", "database_connection"),
		attribute.String("error.action", "retry_recommended"),
	))

	// Set error attributes
	span.SetAttributes(
		attribute.Bool("error", true),
		attribute.String("error.category", "infrastructure"),
		attribute.Int("http.status_code", 500),
		attribute.String("response.status", "error"),
	)

	c.JSON(500, gin.H{
		"error":      "something went wrong",
		"error_code": "DB_CONNECTION_ERROR",
		"request_id": c.GetString("request_id"),
	})
}

func externalAPIHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	ctx, span := tracer.Start(ctx, "externalAPIHandler.processing",
		trace.WithAttributes(
			attribute.String("handler.name", "external_api"),
			attribute.String("handler.type", "aggregator"),
			attribute.Int("api.calls_count", 2),
		),
	)
	defer span.End()

	span.AddEvent("handler.started", trace.WithAttributes(
		attribute.String("state", "initializing"),
	))
	reqLogger.Info().Msg("making external API calls with otelhttp")

	// Create HTTP client with otelhttp instrumentation (auto-creates child spans)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	// === API Call 1: JSONPlaceholder Posts ===
	span.AddEvent("api1.request_started", trace.WithAttributes(
		attribute.String("state", "calling"),
		attribute.String("api", "jsonplaceholder_posts"),
	))
	reqLogger.Info().Msg("calling jsonplaceholder API")

	startTime1 := time.Now()
	req1, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "RequestCreationError"),
			attribute.String("error.api", "api1"),
		))
		span.AddEvent("api1.request_failed", trace.WithAttributes(
			attribute.String("state", "error"),
			attribute.String("error", err.Error()),
		))
		reqLogger.Error().Err(err).Msg("failed to create request")
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	resp1, err := client.Do(req1)
	duration1 := time.Since(startTime1)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "HTTPRequestError"),
			attribute.String("error.api", "api1"),
			attribute.Int64("error.duration_ms", duration1.Milliseconds()),
		))
		span.AddEvent("api1.call_failed", trace.WithAttributes(
			attribute.String("state", "error"),
			attribute.String("error", err.Error()),
		))
		reqLogger.Error().Err(err).Msg("failed to call jsonplaceholder API")
		c.JSON(500, gin.H{"error": "failed to call external API"})
		return
	}
	defer resp1.Body.Close()

	span.AddEvent("api1.response_received", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("status_code", resp1.StatusCode),
		attribute.Int64("duration_ms", duration1.Milliseconds()),
	))
	span.SetAttributes(
		attribute.Int("external.api1.status", resp1.StatusCode),
		attribute.Int64("external.api1.duration_ms", duration1.Milliseconds()),
		attribute.String("external.api1.url", "https://jsonplaceholder.typicode.com/posts/1"),
	)
	reqLogger.Info().Int("status", resp1.StatusCode).Int64("duration_ms", duration1.Milliseconds()).Msg("jsonplaceholder API response")

	// === API Call 2: JSONPlaceholder Users ===
	span.AddEvent("api2.request_started", trace.WithAttributes(
		attribute.String("state", "calling"),
		attribute.String("api", "jsonplaceholder_users"),
	))
	reqLogger.Info().Msg("calling another API")

	startTime2 := time.Now()
	req2, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/users/1", nil)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "RequestCreationError"),
			attribute.String("error.api", "api2"),
		))
		span.AddEvent("api2.request_failed", trace.WithAttributes(
			attribute.String("state", "error"),
			attribute.String("error", err.Error()),
		))
		reqLogger.Error().Err(err).Msg("failed to create request")
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	resp2, err := client.Do(req2)
	duration2 := time.Since(startTime2)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "HTTPRequestError"),
			attribute.String("error.api", "api2"),
			attribute.Int64("error.duration_ms", duration2.Milliseconds()),
		))
		span.AddEvent("api2.call_failed", trace.WithAttributes(
			attribute.String("state", "error"),
			attribute.String("error", err.Error()),
		))
		reqLogger.Error().Err(err).Msg("failed to call users API")
		c.JSON(500, gin.H{"error": "failed to call external API"})
		return
	}
	defer resp2.Body.Close()

	span.AddEvent("api2.response_received", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("status_code", resp2.StatusCode),
		attribute.Int64("duration_ms", duration2.Milliseconds()),
	))
	span.SetAttributes(
		attribute.Int("external.api2.status", resp2.StatusCode),
		attribute.Int64("external.api2.duration_ms", duration2.Milliseconds()),
		attribute.String("external.api2.url", "https://jsonplaceholder.typicode.com/users/1"),
	)
	reqLogger.Info().Int("status", resp2.StatusCode).Int64("duration_ms", duration2.Milliseconds()).Msg("users API response")

	// === Final Summary ===
	totalDuration := duration1 + duration2
	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int("api.successful_calls", 2),
		attribute.Int("api.failed_calls", 0),
		attribute.Int64("api.total_duration_ms", totalDuration.Milliseconds()),
		attribute.String("response.status", "ok"),
	)

	span.AddEvent("handler.completed", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("total_calls", 2),
		attribute.Int64("total_duration_ms", totalDuration.Milliseconds()),
	))

	reqLogger.Info().Msg("external API handler finished")

	c.JSON(200, gin.H{
		"message": "external API calls completed",
		"api1": gin.H{
			"url":         "https://jsonplaceholder.typicode.com/posts/1",
			"status":      resp1.StatusCode,
			"duration_ms": duration1.Milliseconds(),
		},
		"api2": gin.H{
			"url":         "https://jsonplaceholder.typicode.com/users/1",
			"status":      resp2.StatusCode,
			"duration_ms": duration2.Milliseconds(),
		},
		"total_duration_ms": totalDuration.Milliseconds(),
		"request_id":        c.GetString("request_id"),
	})
}

// propagationDemoHandler demonstrates how otelhttp propagates trace context
// to downstream services via HTTP headers (W3C Trace Context standard)
func propagationDemoHandler(c *gin.Context) {
	ctx := c.Request.Context()

	val, exists := c.Get("logger")
	var reqLogger zerolog.Logger
	if exists {
		reqLogger = val.(zerolog.Logger)
	} else {
		reqLogger = log.With().Logger()
	}

	ctx, span := tracer.Start(ctx, "propagationDemoHandler.processing",
		trace.WithAttributes(
			attribute.String("handler.name", "propagation_demo"),
			attribute.String("handler.type", "context_propagation"),
			attribute.String("propagation.standard", "W3C Trace Context"),
		),
	)
	defer span.End()

	span.AddEvent("handler.started", trace.WithAttributes(
		attribute.String("state", "extracting_context"),
	))

	// === 1. Extract current trace context ===
	spanCtx := trace.SpanContextFromContext(ctx)
	traceID := spanCtx.TraceID().String()
	spanID := spanCtx.SpanID().String()
	isSampled := spanCtx.IsSampled()

	span.SetAttributes(
		attribute.String("trace.id", traceID),
		attribute.String("span.id", spanID),
		attribute.Bool("trace.sampled", isSampled),
	)

	span.AddEvent("context.extracted", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.String("trace_id", traceID),
		attribute.String("span_id", spanID),
	))

	reqLogger.Info().
		Str("trace_id", traceID).
		Str("span_id", spanID).
		Bool("sampled", isSampled).
		Msg("extracted trace context from incoming request")

	// === 2. Add custom business metadata to span ===
	businessMetadata := map[string]string{
		"user.id":        "user-12345",
		"session.id":     "sess-abc-789",
		"tenant.id":      "acme-corp",
		"feature.flag":   "new-checkout-enabled",
		"request.source": "mobile-app",
	}

	span.SetAttributes(
		attribute.String("business.user_id", businessMetadata["user.id"]),
		attribute.String("business.session_id", businessMetadata["session.id"]),
		attribute.String("business.tenant_id", businessMetadata["tenant.id"]),
		attribute.String("business.feature_flag", businessMetadata["feature.flag"]),
		attribute.String("business.request_source", businessMetadata["request.source"]),
	)

	span.AddEvent("metadata.attached", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("metadata_count", len(businessMetadata)),
	))

	reqLogger.Info().
		Interface("metadata", businessMetadata).
		Msg("business metadata attached to span")

	// === 3. Create HTTP client with otelhttp (auto-propagates context) ===
	// otelhttp.NewTransport will:
	// - Create child span for each HTTP request
	// - Inject trace context into HTTP headers (traceparent, tracestate)
	// - Propagate context to downstream services
	// - Record HTTP metrics and errors
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	span.AddEvent("http_client.created", trace.WithAttributes(
		attribute.String("state", "ready"),
		attribute.String("transport", "otelhttp"),
	))

	// === 4. Make HTTP request - otelhttp will auto-inject trace headers ===
	span.AddEvent("downstream.call_started", trace.WithAttributes(
		attribute.String("state", "calling"),
		attribute.String("target", "jsonplaceholder.typicode.com"),
	))

	reqLogger.Info().Msg("making downstream HTTP call with context propagation")

	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "RequestCreationError"),
		))
		span.AddEvent("request.creation_failed", trace.WithAttributes(
			attribute.String("state", "error"),
		))
		reqLogger.Error().Err(err).Msg("failed to create HTTP request")
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	// Add custom headers (business context)
	req.Header.Set("X-Request-ID", c.GetString("request_id"))
	req.Header.Set("X-User-ID", businessMetadata["user.id"])
	req.Header.Set("X-Tenant-ID", businessMetadata["tenant.id"])

	// === 5. Execute request - otelhttp injects trace headers automatically ===
	// Headers that will be injected by otelhttp:
	// - traceparent: 00-{trace-id}-{span-id}-{flags}
	// - tracestate: (optional vendor-specific state)
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("error.type", "HTTPRequestError"),
			attribute.Int64("error.duration_ms", duration.Milliseconds()),
		))
		span.AddEvent("downstream.call_failed", trace.WithAttributes(
			attribute.String("state", "error"),
			attribute.String("error", err.Error()),
		))
		reqLogger.Error().Err(err).Msg("downstream HTTP call failed")
		c.JSON(500, gin.H{"error": "downstream call failed"})
		return
	}
	defer resp.Body.Close()

	span.AddEvent("downstream.response_received", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.Int("status_code", resp.StatusCode),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	))

	span.SetAttributes(
		attribute.Int("http.downstream.status_code", resp.StatusCode),
		attribute.Int64("http.downstream.duration_ms", duration.Milliseconds()),
		attribute.String("http.downstream.url", req.URL.String()),
	)

	reqLogger.Info().
		Int("status", resp.StatusCode).
		Int64("duration_ms", duration.Milliseconds()).
		Msg("downstream HTTP call completed")

	// === 6. Demonstrate header injection (for educational purposes) ===
	// Create a dummy request to show what headers would be injected
	dummyReq, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/api", nil)

	// Manually inject trace context (normally done by otelhttp automatically)
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(dummyReq.Header))

	injectedHeaders := make(map[string]string)
	for key, values := range dummyReq.Header {
		if len(values) > 0 {
			injectedHeaders[key] = values[0]
		}
	}

	span.AddEvent("headers.demonstration", trace.WithAttributes(
		attribute.String("state", "completed"),
		attribute.Int("header_count", len(injectedHeaders)),
	))

	reqLogger.Info().
		Interface("injected_headers", injectedHeaders).
		Msg("trace context headers that would be injected")

	// === 7. Create child span to demonstrate context inheritance ===
	ctx, childSpan := tracer.Start(ctx, "propagationDemoHandler.child_operation",
		trace.WithAttributes(
			attribute.String("operation.type", "child_span_demo"),
		),
	)

	childSpanCtx := trace.SpanContextFromContext(ctx)
	childSpan.SetAttributes(
		attribute.String("parent.trace_id", traceID),
		attribute.String("parent.span_id", spanID),
		attribute.String("child.span_id", childSpanCtx.SpanID().String()),
		attribute.Bool("context.inherited", true),
	)

	childSpan.AddEvent("child.processing", trace.WithAttributes(
		attribute.String("state", "working"),
	))

	time.Sleep(50 * time.Millisecond)

	childSpan.AddEvent("child.completed", trace.WithAttributes(
		attribute.String("state", "success"),
	))
	childSpan.End()

	// === 8. Final summary ===
	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.String("propagation.method", "W3C Trace Context"),
		attribute.String("propagation.transport", "otelhttp"),
		attribute.Int("downstream.calls", 1),
		attribute.Int("child.spans", 1),
	)

	span.AddEvent("handler.completed", trace.WithAttributes(
		attribute.String("state", "success"),
		attribute.String("result", "context_propagated"),
	))

	reqLogger.Info().Msg("context propagation demonstration completed")

	// === Response ===
	c.JSON(200, gin.H{
		"message": "Context Propagation Demonstration with otelhttp",
		"trace_context": gin.H{
			"trace_id": traceID,
			"span_id":  spanID,
			"sampled":  isSampled,
			"format":   "W3C Trace Context (traceparent header)",
		},
		"business_metadata": businessMetadata,
		"propagation_mechanism": gin.H{
			"transport":     "otelhttp.NewTransport",
			"standard":      "W3C Trace Context",
			"header_format": "traceparent: 00-{trace-id}-{parent-span-id}-{flags}",
			"automatic":     true,
			"description":   "otelhttp automatically injects trace context headers to all HTTP requests",
		},
		"injected_headers": gin.H{
			"traceparent": injectedHeaders["Traceparent"],
			"tracestate":  injectedHeaders["Tracestate"],
			"explanation": "These headers are automatically added by otelhttp to propagate trace context",
		},
		"downstream_call": gin.H{
			"url":         "https://jsonplaceholder.typicode.com/todos/1",
			"status":      resp.StatusCode,
			"duration_ms": duration.Milliseconds(),
			"context":     "Trace context was automatically propagated via HTTP headers",
		},
		"child_span": gin.H{
			"created":     true,
			"inherited":   "Child span automatically inherits parent trace context",
			"trace_id":    traceID,
			"parent_span": spanID,
			"child_span":  childSpanCtx.SpanID().String(),
		},
		"how_it_works": gin.H{
			"step_1": "otelhttp.NewTransport wraps http.DefaultTransport",
			"step_2": "When making HTTP request, otelhttp creates child span",
			"step_3": "Trace context is extracted from current context",
			"step_4": "Context is injected into HTTP headers (traceparent, tracestate)",
			"step_5": "Downstream service can extract context from headers",
			"step_6": "This creates distributed trace across services",
		},
		"benefits": gin.H{
			"automatic":     "No manual header management required",
			"standardized":  "W3C Trace Context is industry standard",
			"distributed":   "Traces span across multiple services",
			"observability": "Full request flow visible in Jaeger UI",
		},
		"request_id": c.GetString("request_id"),
	})
}
