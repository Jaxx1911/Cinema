package log

import (
	"context"
	"fmt"
)

var globalLogger *logger

func Info(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(addCtxValue(ctx, msg), args...)
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(addCtxValue(ctx, msg), args...)
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(addCtxValue(ctx, msg), args...)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(addCtxValue(ctx, msg), args...)
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(msg, args...)
}

func addCtxValue(ctx context.Context, msg string) string {
	trackId := GetTrackId(ctx)
	if trackId == "" {
		return msg
	}
	return fmt.Sprintf("%s, track_id:[%s]", msg, trackId)
}

func GetTrackId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	trackId := ""
	if ctx.Value("track_id") != nil {
		trackId = ctx.Value("track_id").(string)
	}
	return trackId
}
