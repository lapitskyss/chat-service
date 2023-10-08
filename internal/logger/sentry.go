package logger

import (
	"crypto/tls"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func NewSentryClient(dsn, env, version string) (*sentry.Client, error) {
	return sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     version,
		Environment: env,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	})
}
