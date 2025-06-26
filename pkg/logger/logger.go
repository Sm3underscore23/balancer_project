package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	RequestPath       = "request_path"
	RequestMethod     = "request_method"
	RequestRemoteAddr = "request_remote_address"
	StatusCode        = "status_code"
	Error             = "error"
	HendlerName       = "handler_name"

	ClientID            = "client_id"
	ClientTokenAmount   = "client_token_amount"
	GroupClientLimits   = "client_limis"
	ClinetTokenCapacity = "client_token_capacity"
	ClinetRateRefill    = "client_rate_reffil"

	GroupBackend       = "backend_info"
	BackendUrl         = "backend_url"
	BackendHelthUrl    = "helth_url"
	BackendHelthMethod = "helth_method"
)

type keyType int

const loggerKey = keyType(0)

func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v == nil {
		return slog.Default()
	}

	logger := v.(*slog.Logger)
	return logger
}

func ContextWithSlogLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func AddValuesToContext(ctx context.Context, args ...any) context.Context {
	l := FromContext(ctx)
	l = l.With(args...)
	return ContextWithSlogLogger(ctx, l)
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func InitLogging() {
	handler := slog.Handler(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	))

	slog.SetDefault(slog.New(handler))

}
