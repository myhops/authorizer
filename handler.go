package authorizer

import (
	"fmt"
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
	if a.NotAuthorizedStatusCode == 0 {
		a.NotAuthorizedStatusCode = http.StatusForbidden
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.Authorized(r.Header) {
			w.WriteHeader(a.NotAuthorizedStatusCode)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
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
		keys[h] = append(keys[h],  hk[1])
	}
	for k, v := range keys {
		a.Keys = append(a.Keys, &Key{Header: k, Values: v})
	}
	return a, nil
}
