package authorizer

import (
	"fmt"
	"log/slog"
	"net/http"
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

func (a *Authorizer) Handle() http.HandlerFunc {
	logger := a.Logger.With(slog.String("function", "handler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handler called", slog.String("path", r.URL.RawPath), slog.String("host", r.Host))
		if !a.Authorized(r.Header) {
			logger.Info("forbidden")
			w.WriteHeader(a.NotAuthorizedStatusCode)
			return
		}
		logger.Info("authorized")
		w.WriteHeader(a.AuthorizedStatusCode)
	})
}

func New(headerKeys [][]string) (*Authorizer, error) {
	a := &Authorizer{
		NotAuthorizedStatusCode: http.StatusForbidden,
	}
	keys := map[string][]string{}
	for _, hk := range headerKeys {
		if len(hk) != 2 {
			return nil, fmt.Errorf("expected 2 elements")
		}
		h := http.CanonicalHeaderKey(hk[0])
		keys[h] = append(keys[h], hk[1])
	}
	for k, v := range keys {
		a.Keys = append(a.Keys, &Key{Header: k, Values: v})
	}
	return a, nil
}
