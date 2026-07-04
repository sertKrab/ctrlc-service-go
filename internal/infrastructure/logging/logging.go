// Package logging provides the process-wide structured logger for this
// service. Every layer (handler, usecase, repository) should log through
// slog rather than the stdlib "log" package or a second logging library.
//
// The JSON handler writes to stdout; shipping to a log backend (ELK,
// Loki, etc.) is done by scraping stdout with a sidecar (Filebeat, Fluent
// Bit) rather than the application opening a network client itself. If a
// direct-ship handler is needed later, implement slog.Handler and swap it
// in Init — call sites using slog.InfoContext/WarnContext/ErrorContext do
// not change.
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type ctxKey struct{}

// Standard structured-log field names. Use these constants instead of
// literal strings so field names stay identical across modules/subagent
// runs (e.g. don't let one module log "user_id" and another "userId").
const (
	// KeyFunc identifies which function/method emitted the log line, e.g.
	// "LoginUseCase.Execute". Declare once per function:
	//   log := logging.FromContext(ctx).With(logging.KeyFunc, "LoginUseCase.Execute")
	// then reuse `log` for every log call in that function.
	KeyFunc = "func"

	// KeyUserID identifies which user/actor the log line is about.
	KeyUserID = "user_id"
)

// Init sets the process-wide default logger to a JSON handler on stdout at
// the given level ("debug", "info", "warn", "error" — defaults to info on
// an unrecognized value) and returns it.
func Init(level string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	}))
	slog.SetDefault(logger)
	return logger
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithRequestID attaches a request ID to ctx so FromContext can include it
// in every subsequent log call without repeating "request_id" at each call
// site. Call this once per request (the RequestID middleware does this).
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, requestID)
}

// FromContext returns the default logger, pre-annotated with the request ID
// stored in ctx if RequestID middleware set one. Safe to call with any ctx —
// falls back to the plain default logger when no request ID is present.
func FromContext(ctx context.Context) *slog.Logger {
	id, _ := ctx.Value(ctxKey{}).(string)
	if id == "" {
		return slog.Default()
	}
	return slog.Default().With("request_id", id)
}

// MaskOptions controls how Mask redacts a sensitive value.
type MaskOptions struct {
	KeepPrefix int  // characters kept visible at the start
	KeepSuffix int  // characters kept visible at the end
	MaskChar   rune // character used to replace the redacted middle
}

// DefaultMaskOptions is used when Mask is called with no explicit options —
// a generic redaction for fields the SRS/testcase spec does not call out a
// specific masking format for.
var DefaultMaskOptions = MaskOptions{KeepPrefix: 2, KeepSuffix: 2, MaskChar: '*'}

// Mask redacts s for logging, keeping only the configured prefix/suffix
// visible. Pass explicit MaskOptions when the SRS specifies a field's
// masking format (e.g. "show only the last 4 digits of the card number" →
// logging.Mask(cardNo, logging.MaskOptions{KeepSuffix: 4})); omit opts for
// fields with no specified format — DefaultMaskOptions applies.
// Never log a sensitive value (citizen ID, card number, phone, full email,
// token, password) without passing it through Mask first.
func Mask(s string, opts ...MaskOptions) string {
	o := DefaultMaskOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.MaskChar == 0 {
		o.MaskChar = '*'
	}
	runes := []rune(s)
	n := len(runes)
	if n <= o.KeepPrefix+o.KeepSuffix {
		return strings.Repeat(string(o.MaskChar), n)
	}
	masked := make([]rune, n)
	copy(masked, runes[:o.KeepPrefix])
	for i := o.KeepPrefix; i < n-o.KeepSuffix; i++ {
		masked[i] = o.MaskChar
	}
	copy(masked[n-o.KeepSuffix:], runes[n-o.KeepSuffix:])
	return string(masked)
}
