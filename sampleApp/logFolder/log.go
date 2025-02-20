package logFolder

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"os"
	"sample_app/config"
	"strconv"
)

// Logger enable you to log actions
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningln(args ...interface{})
	Warnln(args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

type loggerWrapper struct {
	*logrus.Logger
}

func (lw loggerWrapper) WithFields(fields map[string]interface{}) Logger {
	return loggerWrapper{Logger: lw.Logger.WithFields(fields).Logger}
}

var defaultLogger loggerWrapper
var env string

func init() {
	cfg := config.Config()
	defaultLogger = loggerWrapper{Logger: newLogrusLogger(cfg)}
	env = "stg"
	if cfg.GetString("ENV") == "prod" {
		env = "prd"
	}
}

func newLogrusLogger(cfg config.Provider) *logrus.Logger {
	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.99-07:00",
	})

	l.SetOutput(os.Stdout)

	if cfg.GetBool("JSON_LOGS") {
		l.SetFormatter(&logrus.JSONFormatter{})
	}
	l.Out = os.Stderr

	switch cfg.GetString("LOG_LEVEL") {
	case "debug":
		l.Level = logrus.DebugLevel
	case "warning":
		l.Level = logrus.WarnLevel
	case "info":
		l.Level = logrus.InfoLevel
	case "error":
		l.Level = logrus.ErrorLevel
	case "panic":
		l.Level = logrus.PanicLevel
	default:
		l.Level = logrus.InfoLevel
	}

	return l
}

// Debug package-level convenience method.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf package-level convenience method.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Debugln package-level convenience method.
func Debugln(args ...interface{}) {
	defaultLogger.Debugln(args...)
}

// Error package-level convenience method.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf package-level convenience method.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Errorln package-level convenience method.
func Errorln(args ...interface{}) {
	defaultLogger.Errorln(args...)
}

// Fatal package-level convenience method.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf package-level convenience method.
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Fatalln package-level convenience method.
func Fatalln(args ...interface{}) {
	defaultLogger.Fatalln(args...)
}

// Info package-level convenience method.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof package-level convenience method.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Infoln package-level convenience method.
func Infoln(args ...interface{}) {
	defaultLogger.Infoln(args...)
}

// Panic package-level convenience method.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Panicf package-level convenience method.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

// Panicln package-level convenience method.
func Panicln(args ...interface{}) {
	defaultLogger.Panicln(args...)
}

// Print package-level convenience method.
func Print(args ...interface{}) {
	defaultLogger.Print(args...)
}

// Printf package-level convenience method.
func Printf(format string, args ...interface{}) {
	defaultLogger.Printf(format, args...)
}

// Println package-level convenience method.
func Println(args ...interface{}) {
	defaultLogger.Println(args...)
}

// Warn package-level convenience method.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf package-level convenience method.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Warning package-level convenience method.
func Warning(args ...interface{}) {
	defaultLogger.Warning(args...)
}

// Warningf package-level convenience method.
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Warningln package-level convenience method.
func Warningln(args ...interface{}) {
	defaultLogger.Warningln(args...)
}

// Warnln package-level convenience method.
func Warnln(args ...interface{}) {
	defaultLogger.Warnln(args...)
}

// WithFields package-level convenience method.
func WithFields(fields map[string]interface{}, args ...interface{}) Logger {
	return defaultLogger.WithFields(fields)
}

// StandardLogFields creates log fields for a given context and span.
func StandardLogFields(span trace.Span, serviceName string) map[string]interface{} {
	if span == nil {
		return logrus.Fields{}
	}
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return logrus.Fields{
		"dd.trace_id": convertTraceID(traceID),
		"dd.span_id":  convertTraceID(spanID),
		"dd.service":  serviceName,
		"dd.env":      env,
	}
}

// Formatting for Datadog
func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
