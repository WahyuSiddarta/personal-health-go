package middleware

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// CaptureException captures an exception in Sentry with request context
func CaptureException(c echo.Context, err error) {
	if err == nil {
		return
	}

	hub := sentry.GetHubFromContext(c.Request().Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err)
}

// CaptureMessage captures a message in Sentry with request context
func CaptureMessage(c echo.Context, message string) {
	hub := sentry.GetHubFromContext(c.Request().Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureMessage(message)
}

// CaptureError captures an error in Sentry with additional context
func CaptureError(c echo.Context, err error, tags map[string]string, extra map[string]interface{}) {
	if err == nil {
		return
	}

	hub := sentry.GetHubFromContext(c.Request().Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	// Add tags
	for key, value := range tags {
		hub.Scope().SetTag(key, value)
	}

	// Add extra context
	for key, value := range extra {
		hub.Scope().SetExtra(key, value)
	}

	// Add request information
	hub.Scope().SetContext("request", map[string]interface{}{
		"method": c.Request().Method,
		"url":    c.Request().URL.String(),
		"path":   c.Request().URL.Path,
		"query":  c.Request().URL.RawQuery,
	})

	hub.CaptureException(err)
}

// CaptureRecovery captures a panic recovery in Sentry
func CaptureRecovery(c echo.Context, err error, stackTrace string) {
	if err == nil {
		return
	}

	hub := sentry.GetHubFromContext(c.Request().Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.Scope().SetContext("panic", map[string]interface{}{
		"type":        fmt.Sprintf("%T", err),
		"stack_trace": stackTrace,
	})

	hub.Scope().SetLevel(sentry.LevelFatal)

	hub.CaptureException(err)
}

// SetUserContext sets user information in Sentry
func SetUserContext(c echo.Context, userID interface{}, userEmail string) {
	hub := sentry.GetHubFromContext(c.Request().Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.Scope().SetUser(sentry.User{
		ID:    fmt.Sprintf("%v", userID),
		Email: userEmail,
	})
}

// FlushSentry flushes pending events to Sentry
// Call this during graceful shutdown
func FlushSentry(timeoutSeconds int) bool {
	return sentry.Flush(time.Duration(timeoutSeconds) * time.Second)
}
