package grpc

import (
	"context"
	"time"

	"github.com/imrenagicom/demo-app/internal/instrumentation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerTracingInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if tracer == nil {
			return handler(ctx, req)
		}

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.service", "course-service"),
			),
		)
		defer span.End()

		instrumentation.AddEvent(span, "grpc.request.started", instrumentation.StateStarted)

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		span.SetAttributes(
			attribute.Int64("rpc.duration_ms", duration.Milliseconds()),
		)

		if err != nil {
			st, _ := status.FromError(err)
			instrumentation.RecordError(span, err, "GRPCError", "grpc_handler",
				attribute.String("rpc.method", info.FullMethod),
				attribute.Int("rpc.grpc.status_code", int(st.Code())),
			)
			instrumentation.AddEvent(span, "grpc.request.failed", instrumentation.StateFailed,
				attribute.String("error", err.Error()),
				attribute.String("grpc.code", st.Code().String()),
			)
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
		} else {
			instrumentation.AddEvent(span, "grpc.request.completed", instrumentation.StateSuccess)
			span.SetAttributes(
				attribute.Int("rpc.grpc.status_code", int(codes.OK)),
				attribute.Bool("success", true),
			)
		}

		return resp, err
	}
}

func StreamServerTracingInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if tracer == nil {
			return handler(srv, ss)
		}

		ctx, span := tracer.Start(ss.Context(), info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.service", "course-service"),
				attribute.Bool("rpc.is_stream", true),
			),
		)
		defer span.End()

		instrumentation.AddEvent(span, "grpc.stream.started", instrumentation.StateStarted)

		start := time.Now()
		err := handler(srv, &tracedServerStream{ServerStream: ss, ctx: ctx})
		duration := time.Since(start)

		span.SetAttributes(
			attribute.Int64("rpc.duration_ms", duration.Milliseconds()),
		)

		if err != nil {
			st, _ := status.FromError(err)
			instrumentation.RecordError(span, err, "GRPCStreamError", "grpc_stream_handler",
				attribute.String("rpc.method", info.FullMethod),
			)
			instrumentation.AddEvent(span, "grpc.stream.failed", instrumentation.StateFailed)
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
		} else {
			instrumentation.AddEvent(span, "grpc.stream.completed", instrumentation.StateSuccess)
			span.SetAttributes(
				attribute.Int("rpc.grpc.status_code", int(codes.OK)),
				attribute.Bool("success", true),
			)
		}

		return err
	}
}

type tracedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *tracedServerStream) Context() context.Context {
	return s.ctx
}
