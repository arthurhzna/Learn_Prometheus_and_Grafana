package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "status_code"},
	)

	grpcRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_request_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"service", "method", "status_code"},
	)

	grpcRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_server_requests_in_flight",
			Help: "Number of gRPC requests currently being processed",
		},
		[]string{"service", "method"},
	)

	bookingStatusTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_status_total",
			Help: "Total number of bookings by status and course",
		},
		[]string{"course_id", "course_name", "status"},
	)

	bookingReservedGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "booking_reserved_current",
			Help: "Current number of reserved bookings per course",
		},
		[]string{"course_id", "course_name"},
	)

	bookingExpiredGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "booking_expired_current",
			Help: "Current number of expired bookings per course",
		},
		[]string{"course_id", "course_name"},
	)

	bookingPotentialSales = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "booking_potential_sales_amount",
			Help: "Potential sales amount from successful reservations per course",
		},
		[]string{"course_id", "course_name", "currency"},
	)

	bookingCompletedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_completed_total",
			Help: "Total number of completed bookings",
		},
		[]string{"course_id", "course_name"},
	)

	bookingCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_created_total",
			Help: "Total number of created bookings",
		},
		[]string{"course_id", "course_name"},
	)
)

func init() {
	exporter, err := promexporter.New()
	if err != nil {
		panic(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		service, method := splitMethodName(info.FullMethod)

		grpcRequestsInFlight.WithLabelValues(service, method).Inc()
		defer grpcRequestsInFlight.WithLabelValues(service, method).Dec()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		statusCode := status.Code(err).String()

		grpcRequestDuration.WithLabelValues(service, method, statusCode).Observe(duration)
		grpcRequestTotal.WithLabelValues(service, method, statusCode).Inc()

		return resp, err
	}
}

func splitMethodName(fullMethod string) (string, string) {
	var service, method string

	if len(fullMethod) > 0 && fullMethod[0] == '/' {
		fullMethod = fullMethod[1:]
	}

	lastSlash := -1
	for i, c := range fullMethod {
		if c == '/' {
			lastSlash = i
			break
		}
	}

	if lastSlash >= 0 {
		servicePart := fullMethod[:lastSlash]
		method = fullMethod[lastSlash+1:]

		lastDot := -1
		for i := len(servicePart) - 1; i >= 0; i-- {
			if servicePart[i] == '.' {
				lastDot = i
				break
			}
		}
		if lastDot >= 0 {
			service = servicePart[lastDot+1:]
		} else {
			service = servicePart
		}
	} else {
		service = "unknown"
		method = fullMethod
	}

	return service, method
}

func RecordBookingCreated(courseID, courseName string) {
	bookingCreatedTotal.WithLabelValues(courseID, courseName).Inc()
	bookingStatusTotal.WithLabelValues(courseID, courseName, "created").Inc()
}

// RecordBookingReserved records a booking reservation event
func RecordBookingReserved(courseID, courseName string, price float64, currency string) {
	bookingStatusTotal.WithLabelValues(courseID, courseName, "reserved").Inc()
	bookingReservedGauge.WithLabelValues(courseID, courseName).Inc()
	bookingPotentialSales.WithLabelValues(courseID, courseName, currency).Add(price)
}

// RecordBookingExpired records a booking expiration event
func RecordBookingExpired(courseID, courseName string, price float64, currency string) {
	bookingStatusTotal.WithLabelValues(courseID, courseName, "expired").Inc()
	bookingExpiredGauge.WithLabelValues(courseID, courseName).Inc()

	// Decrease reserved count and potential sales
	bookingReservedGauge.WithLabelValues(courseID, courseName).Dec()
	bookingPotentialSales.WithLabelValues(courseID, courseName, currency).Sub(price)
}

// RecordBookingCompleted records a booking completion event
func RecordBookingCompleted(courseID, courseName string, price float64, currency string) {
	bookingStatusTotal.WithLabelValues(courseID, courseName, "completed").Inc()
	bookingCompletedTotal.WithLabelValues(courseID, courseName).Inc()

	// Decrease reserved count (assuming it was reserved before completed)
	bookingReservedGauge.WithLabelValues(courseID, courseName).Dec()
	bookingPotentialSales.WithLabelValues(courseID, courseName, currency).Sub(price)
}

// RecordBookingFailed records a booking failure event
func RecordBookingFailed(courseID, courseName string) {
	bookingStatusTotal.WithLabelValues(courseID, courseName, "failed").Inc()
}
