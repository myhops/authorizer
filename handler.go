package authorizer

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Key contains the valid keys and the name of the header for the key.
type Key struct {
	Header string
	Values []string
}

// InHeader returns true if the http header contains a valid key.
func (k *Key) InHeader(header http.Header) bool {
	values := header[k.Header]
	if values == nil {
		return false
	}
	for _, value := range values {
		for _, v := range k.Values {
			if value == v {
				return true
			}
		}
	}
	return false
}

// Authorizer is a handler that checks requests for valid keys.
type Authorizer struct {
	// Keys contains the list of valid keys.
	Keys                    []*Key
	// NotAuthorizedStatusCode contains the status code the hander 
	// returns when the request does not contain a valid key.
	NotAuthorizedStatusCode int
	// AuthorizedStatusCode contains the status code the hander 
	// returns when the request does not contain a valid key.
	AuthorizedStatusCode    int
	// Logger is the structured logger for handler.
	Logger                  *slog.Logger
}

// Authorized returns true if the http header contains a valid key.
func (a *Authorizer) Authorized(h http.Header) bool {
	for _, k := range a.Keys {
		if k.InHeader(h) {
			return true
		}
	}
	return false
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

// Handle returns a handler func that checks if the request 
// contains a valid api key.
// It returns AuthorizedStatusCode if this is the case and
// NotAuthorizedStatusCode if not.
func (a *Authorizer) Handle() http.HandlerFunc {
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
		if !a.Authorized(r.Header) {
			logger.Info("forbidden")
			w.WriteHeader(a.NotAuthorizedStatusCode)
			return
		}
		logger.Info("authorized")
		w.WriteHeader(a.AuthorizedStatusCode)
	})
}

// AggrateKeys groups the keys with similar headers in one key.
// It transforms the header to the canonical format.
func AggrateKeys(headerKeys [][]string) ([]*Key, error) {
	// Aggregate the keys per header
	keys := map[string][]string{}
	for _, hk := range headerKeys {
		if len(hk) != 2 {
			return nil, fmt.Errorf("expected 2 elements")
		}
		h := http.CanonicalHeaderKey(hk[0])
		keys[h] = append(keys[h], hk[1])
	}
	var res []*Key
	for k, v := range keys {
		res = append(res, &Key{Header: k, Values: v})
	}
	return res, nil
}

// New returns an authorizer with the give keys.
func New(keys []*Key) (*Authorizer, error) {
	a := &Authorizer{
		NotAuthorizedStatusCode: http.StatusForbidden,
		AuthorizedStatusCode:    http.StatusOK,
	}

	a.Keys = keys
	return a, nil
}
