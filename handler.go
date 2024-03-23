package authorizer

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Key struct {
	Header string
	Values []string
}

type Keys struct {
}

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

type Authorizer struct {
	Keys                    []*Key
	NotAuthorizedStatusCode int
	AuthorizedStatusCode    int
	Logger                  *slog.Logger
}

func (a *Authorizer) Authorized(h http.Header) bool {
	for _, k := range a.Keys {
		if k.InHeader(h) {
			return true
		}
	}
	return false
}

// slog valuer for duration as string
type durationLV time.Duration

func (d durationLV) String() string {
	return time.Duration(d).String()
}

func (d durationLV) LogValue() slog.Value {
	return slog.StringValue(d.String())
}


func (a *Authorizer) Handle() http.HandlerFunc {
	logger := a.Logger.With(slog.String("function", "handler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			since := time.Since(start)
			logger.Info("handler called",
				slog.String("path", r.URL.RawPath),
				slog.String("host", r.Host),
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

func New(keys []*Key) (*Authorizer, error) {
	a := &Authorizer{
		NotAuthorizedStatusCode: http.StatusForbidden,
		AuthorizedStatusCode:    http.StatusOK,
	}

	a.Keys = keys
	return a, nil
}
