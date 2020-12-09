package errors

import (
	"log"
	"os"
	"time"

	sentry "github.com/getsentry/sentry-go"
)

var sentryEnabled = false

func initSentry() {
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			log.Fatalf("Failed sentry.Init: %+v", err)
		}
		sentryEnabled = true
	}
}

func EmitSentry(err error) {
	if sentryEnabled {
		eventID := sentry.CaptureException(err)
		if eventID != nil {
			// Add sentry eventID to original error
			if e, ok := err.(*Error); ok {
				_ = e.With("sentry.eventID", eventID)
			}
		}
	}
}

func FlushSentry() {
	if sentryEnabled {
		sentry.Flush(2 * time.Second)
	}
}
