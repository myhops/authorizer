package authorizer

import (
	"log/slog"
	"net/http"
	"time"
)

// Authorizer is a handler that checks requests for valid keys.
type Authorizer struct {
	// NotAuthorizedStatusCode contains the status code the hander
	// returns when the request does not contain a valid key.
	NotAuthorizedStatusCode int
	// AuthorizedStatusCode contains the status code the hander
	// returns when the request does not contain a valid key.
	AuthorizedStatusCode int
	// Logger is the structured logger for handler.
	Logger *slog.Logger
}

// durationLV is a slog valuer for duration as string.
type durationLV time.Duration

// String returns the value of d as fmt formatted string.
func (d durationLV) String() string {
	return time.Duration(d).String()
}

// LogValue returns the value of d.
func (d durationLV) LogValue() slog.Value {
	return slog.StringValue(d.String())
}

type RequestAuthorizer interface {
	Authorize(*http.Request) error
}

// Handle returns a handler func that checks if the request
// contains a valid api key.
// It returns AuthorizedStatusCode if this is the case and
// NotAuthorizedStatusCode if not.
func (a *Authorizer) Handle(ra RequestAuthorizer) http.HandlerFunc {
	logger := a.Logger.With(slog.String("function", "handler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			slog.String("path", r.URL.RawPath),
			slog.String("host", r.Host))

		defer func(start time.Time) {
			since := time.Since(start)
			logger.Info("handler called",
				slog.Time("start", start),
				slog.Duration("duration", since),
				slog.Attr{Key: "duration_string", Value: durationLV(since).LogValue()},
			)
		}(time.Now())
		if err := ra.Authorize(r); err != nil {
			logger.Info("forbidden")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(a.NotAuthorizedStatusCode)
			w.Write([]byte(err.Error()))
			return
		}
		logger.Info("authorized")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(a.AuthorizedStatusCode)
		w.Write([]byte("Authorized"))
	})
}

// New returns an authorizer with the give keys.
func New() (*Authorizer, error) {
	a := &Authorizer{

		NotAuthorizedStatusCode: http.StatusForbidden,
		AuthorizedStatusCode:    http.StatusOK,
	}
	return a, nil
}
