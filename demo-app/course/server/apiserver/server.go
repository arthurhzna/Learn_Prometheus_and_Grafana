package apiserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/imrenagicom/demo-app/course/booking"
	"github.com/imrenagicom/demo-app/course/catalog"
	bookingsrv "github.com/imrenagicom/demo-app/course/server/booking"
	catalogsrv "github.com/imrenagicom/demo-app/course/server/catalog"
	"github.com/imrenagicom/demo-app/internal/config"
	grpcutil "github.com/imrenagicom/demo-app/internal/grpc"
	"github.com/imrenagicom/demo-app/internal/instrumentation"
	"github.com/imrenagicom/demo-app/internal/util"
	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/imrenagicom/demo-app/internal/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var serviceTelemetryName = "course-service"

type ServerOpts struct {
	Clients      *util.Clients
	Config       config.Server
	OTLPEndpoint string
	Tracer       trace.Tracer
}

func NewServer(opts ServerOpts) Server {
	log.Debug().
		Str("postgres", fmt.Sprintf("%s:%s/%s", opts.Config.DB.Host, opts.Config.DB.Port, opts.Config.DB.Name)).
		Str("redis", opts.Config.Redis.Addr()).
		Msg("checking config")

	// s := Server{
	// 	opts:    opts,
	// 	clients: opts.Clients,
	// }

	s := Server{
		opts:                 opts,
		clients:              opts.Clients,
		otlpCollectorAddress: opts.OTLPEndpoint,
		tracer:               opts.Tracer, // ← HARUS ADA INI!
	}

	s.catalogStore = catalog.NewStore(opts.Clients.DB, opts.Clients.Redis)
	s.catalogService = catalog.NewService(s.catalogStore, opts.Clients.DB)
	s.bookingStore = booking.NewStore(opts.Clients.DB, opts.Clients.Redis)
	s.bookingService = booking.NewService(
		opts.Clients.DB,
		s.bookingStore,
		s.catalogStore,
	)
	return s
}

type Server struct {
	opts                 ServerOpts
	clients              *util.Clients
	otlpCollectorAddress string
	tracer               trace.Tracer

	bookingService *booking.Service
	bookingStore   *booking.Store
	catalogService *catalog.Service
	catalogStore   *catalog.Store
}

func (s *Server) Run(ctx context.Context) error {
	log.Info().Msg("starting server")

	if s.tracer != nil {
		log.Info().
			Str("otlp_endpoint", s.otlpCollectorAddress).
			Msg("OpenTelemetry tracing enabled")
	}

	grpcServer := s.newGRPCServer(ctx)
	go func() {
		log.Info().Msgf("initializing grpc server on %s", s.opts.Config.GRPC.Addr())
		lis, err := net.Listen("tcp", s.opts.Config.GRPC.Addr())
		if err != nil {
			log.Fatal().Msgf("failed to listen: %v", err)
		}
		log.Info().Msgf("starting grpc server on %s", s.opts.Config.GRPC.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("unable to start grpc server")
		}
	}()

	httpServer := s.newHTTPServer(ctx)
	go func() {
		log.Info().Msgf("Starting http server for serving gRPC-Gateway and OpenAPI Documentation on %s", s.opts.Config.HTTP.Addr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	gracefulShutdownPeriod := 30 * time.Second

	log.Warn().Msg("shutting down http server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownPeriod)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}
	log.Warn().Msg("http server gracefully stopped")

	log.Warn().Msg("shutting down grpc server")
	grpcServer.GracefulStop()
	log.Warn().Msg("grpc server gracefully stopped")

	log.Warn().Msg("clean up storage")
	if err := s.catalogStore.Clear(); err != nil {
		log.Warn().Err(err).Msg("failed to clear concert store")
	}
	if err := s.bookingStore.Clear(); err != nil {
		log.Warn().Err(err).Msg("failed to clear concert store")
	}
	return nil
}

func (s *Server) newGRPCServer(ctx context.Context) *grpc.Server {

	interceptors := []grpc.UnaryServerInterceptor{
		grpcutil.UnaryServerAppLoggerInterceptor(),
		grpcutil.UnaryServerGRPCLoggerInterceptor(),
		metrics.UnaryServerInterceptor(),
	}

	if s.tracer != nil {
		interceptors = append(interceptors, grpcutil.UnaryServerTracingInterceptor(s.tracer))
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpcutil.StreamServerAppLoggerInterceptor(),
		grpcutil.StreamServerGRPCLoggerInterceptor(),
	}

	if s.tracer != nil {
		streamInterceptors = append(streamInterceptors, grpcutil.StreamServerTracingInterceptor(s.tracer))
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}

	grpcServer := grpc.NewServer(opts...)
	bookingSrv := bookingsrv.New(s.bookingService)
	catalogSrv := catalogsrv.New(s.catalogService)
	v1.RegisterBookingServiceServer(grpcServer, bookingSrv)
	v1.RegisterCatalogServiceServer(grpcServer, catalogSrv)
	return grpcServer
}

func (s *Server) newHTTPServer(ctx context.Context) *http.Server {
	gRPCEndpoint := s.opts.Config.GRPC.Addr()
	conn, err := grpc.DialContext(
		ctx,
		gRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to dial grpc server: %v", err)
	}

	gwmux := runtime.NewServeMux()
	mustRegisterGWHandler(ctx, v1.RegisterCatalogServiceHandler, gwmux, conn)
	mustRegisterGWHandler(ctx, v1.RegisterBookingServiceHandler, gwmux, conn)

	mux := mux.NewRouter()
	mux.HandleFunc("/healthz", s.healthz())
	mux.HandleFunc("/readyz", s.readyz())

	mux.Handle("/metrics", promhttp.Handler())

	mux.PathPrefix("/debug/").Handler(http.DefaultServeMux)

	api := mux.PathPrefix("/api/course").Subrouter()

	if s.tracer != nil {
		api.Use(s.httpTracingMiddleware)
	}

	api.PathPrefix("/v1").Handler(gwmux)

	sh := http.StripPrefix("/swagger/",
		http.FileServer(http.Dir("./third_party/OpenAPI/")))
	mux.PathPrefix("/swagger/").Handler(sh)

	gwServer := &http.Server{
		Addr:    s.opts.Config.HTTP.Addr(),
		Handler: mux,
	}
	return gwServer
}

func (s *Server) httpTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.tracer == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx, span := s.tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.target", r.URL.Path),
				attribute.String("http.client_ip", r.RemoteAddr),
			),
		)
		defer span.End()

		instrumentation.AddEvent(span, "http.request.started", instrumentation.StateStarted)

		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		duration := time.Since(start)

		span.SetAttributes(
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
			attribute.Bool("success", true),
		)

		instrumentation.AddEvent(span, "http.request.completed", instrumentation.StateSuccess,
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)
	})
}

type registerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

func mustRegisterGWHandler(ctx context.Context, register registerFunc, mux *runtime.ServeMux, conn *grpc.ClientConn) {
	err := register(ctx, mux, conn)
	if err != nil {
		panic(err)
	}
}

func (s *Server) healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) readyz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
