package authorizer

import (
	"errors"
	"fmt"
	"net/http"
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

type Keys []*Key

var (
	ErrNotAuthorized = errors.New("not authorized")
)

func (k Keys) authorizeHeader(h http.Header) error {
	for _, kk := range k {
		if kk.InHeader(h) {
			return nil
		}
	}
	return ErrNotAuthorized
}

// Authorized returns true if the http header contains a valid key.
func (k Keys) Authorize(r *http.Request) error {
	return k.authorizeHeader(r.Header)
}

// NewFromKeyList groups the keys with similar headers in one key.
// It transforms the header to the canonical format.
func NewFromKeyList(kl [][]string) (Keys, error) {
	// Aggregate the keys per header
	keys := map[string][]string{}
	for _, hk := range kl {
		if len(hk) != 2 {
			return nil, fmt.Errorf("expected 2 elements")
		}
		h := http.CanonicalHeaderKey(hk[0])
		keys[h] = append(keys[h], hk[1])
	}
	var res Keys
	for k, v := range keys {
		res = append(res, &Key{Header: k, Values: v})
	}
	return res, nil
}
