package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/requestContext"
)

type Logger struct {
	base *slog.Logger
}

func New() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return &Logger{
		base: slog.New(handler),
	}
}

func (l *Logger) withContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return l.base
	}

	rc, ok := ctx.Value(requestContext.RequestKey).(requestContext.ReqContext)
	if !ok {
		return l.base
	}

	args := []any{
		"req_id", rc.ReqID,
	}

	if rc.UserID != nil {
		args = append(args, "user_id", rc.UserID)
	}

	if rc.UserEmail != nil {
		args = append(args, "email", rc.UserEmail)
	}

	return l.base.With(args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Info(msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Error(msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Debug(msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Warn(msg, args...)
}

func (l *Logger) HTTP(ctx context.Context, method string, path string) {
	l.withContext(ctx).Info("http_request",
		"method", method,
		"path", path,
	)
}

func (l *Logger) Auth(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Info("auth_"+msg, args...)
}

func (l *Logger) Event(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Info("event_"+msg, args...)
}
