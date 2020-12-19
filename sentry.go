package golambda

import (
	"fmt"
	"log"
	"os"
	"time"

	sentry "github.com/getsentry/sentry-go"
)

var sentryEnabled = false

func initSentry() {
	if dsn, ok := os.LookupEnv("SENTRY_DSN"); ok {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			log.Fatalf("Failed sentry.Init: %+v", err)
		}
		sentryEnabled = true
	}
}

func emitSentry(err error) string {
	if sentryEnabled {
		eventID := sentry.CaptureException(err)
		if eventID != nil {
			// Add sentry eventID to original error
			return fmt.Sprintf("%v", *eventID)
		}
	}
	return ""
}

func flushSentry() {
	if sentryEnabled {
		sentry.Flush(2 * time.Second)
	}
}
