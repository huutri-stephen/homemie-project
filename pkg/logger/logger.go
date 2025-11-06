package logger

import (
	"context"
	"go.uber.org/zap"
)

type contextKey string

const (
	SESSION_LOGGER contextKey = "SESSION_LOGGER"
	TRACE_ID       contextKey = "TRACE_ID"
)

var (
	loggerInstance *zap.Logger
)

func GetSessionLogger(fields ...zap.Field) *zap.SugaredLogger {
	return loggerInstance.With(fields...).Sugar()
}

// getSessionLoggerFromCtx returns a *zap.Logger from the context, or the global logger.
func getSessionLoggerFromCtx(ctx context.Context) *zap.Logger {
	sessionLogger, ok := ctx.Value(SESSION_LOGGER).(*zap.Logger) // Change to *zap.Logger
	if ok {
		return sessionLogger
	}

	logger := loggerInstance
	if traceID, ok := ctx.Value(TRACE_ID).(string); ok && traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}
	return logger
}

func Info(ctx context.Context, msg string) {
	getSessionLoggerFromCtx(ctx).Sugar().Info(msg)
}
func Infof(ctx context.Context, format string, args ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Infof(format, args...)
}
func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Infow(msg, keysAndValues...)
}

func Warn(ctx context.Context, msg string) {
	getSessionLoggerFromCtx(ctx).Sugar().Warn(msg)
}
func Warnf(ctx context.Context, format string, args ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Warnf(format, args...)
}
func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Warnw(msg, keysAndValues...)
}

func Error(ctx context.Context, msg string) {
	getSessionLoggerFromCtx(ctx).Sugar().Error(msg)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Errorf(format, args...)
}
func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Errorw(msg, keysAndValues...)
}

func Debug(ctx context.Context, msg string) {
	getSessionLoggerFromCtx(ctx).Sugar().Debug(msg)
}
func Debugf(ctx context.Context, format string, args ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Debugf(format, args...)
}
func Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	getSessionLoggerFromCtx(ctx).Sugar().Debugw(msg, keysAndValues...)
}

func NewLogContextFromParent(parent context.Context) (context.Context, context.CancelFunc) {
	newCtx := context.Background()
	if parentLogger := parent.Value(SESSION_LOGGER); parentLogger != nil {
		newCtx = context.WithValue(newCtx, SESSION_LOGGER, parentLogger)
	}
	if parentTraceID := parent.Value(TRACE_ID); parentTraceID != nil {
		newCtx = context.WithValue(newCtx, TRACE_ID, parentTraceID)
	}
	return context.WithCancel(newCtx)
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TRACE_ID, traceID)
}

func ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	currentLogger := getSessionLoggerFromCtx(ctx)
	newLogger := currentLogger.With(fields...)
	return context.WithValue(ctx, SESSION_LOGGER, newLogger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	return getSessionLoggerFromCtx(ctx).Sugar()
}