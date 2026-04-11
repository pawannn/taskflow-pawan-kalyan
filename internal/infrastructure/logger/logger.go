package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/requestContext"
)

type Logger struct {
	http  *slog.Logger
	auth  *slog.Logger
	event *slog.Logger
	error *slog.Logger
}

func newFileLogger(filename string) *slog.Logger {
	_ = os.MkdirAll("logs", os.ModePerm)

	file, err := os.OpenFile("logs/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(handler)
}

func New() *Logger {
	return &Logger{
		http:  newFileLogger("http.log"),
		auth:  newFileLogger("auth.log"),
		event: newFileLogger("event.log"),
		error: newFileLogger("error.log"),
	}
}

func (l *Logger) withContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if l == nil || logger == nil {
		return slog.Default()
	}

	if ctx == nil {
		return logger
	}

	rc, ok := ctx.Value(requestContext.RequestKey).(requestContext.ReqContext)
	if !ok {
		return logger
	}

	args := []any{
		"req_id", rc.ReqID,
	}

	if rc.UserID != "" {
		args = append(args, "user_id", rc.UserID)
	}

	if rc.UserEmail != "" {
		args = append(args, "email", rc.UserEmail)
	}

	return logger.With(args...)
}

func (l *Logger) HTTP(ctx context.Context, method string, path string) {
	l.withContext(ctx, l.http).Info("http_request",
		"method", method,
		"path", path,
	)
}

func (l *Logger) Auth(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.auth).Info("auth_"+msg, args...)
}

func (l *Logger) Event(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.event).Info("event_"+msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.error).Error(msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.event).Debug(msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.event).Warn(msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx, l.event).Info(msg, args...)
}
