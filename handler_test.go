package authorizer

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	ra, err := NewFromKeyList([][]string{{"x-api-key", "valid-key"}})
	if err != nil {
		t.Fatalf("error creating request authorizer: %s", err.Error())
	}

	a := &Authorizer{
		Logger:                  slog.Default(),
		NotAuthorizedStatusCode: http.StatusForbidden,
		AuthorizedStatusCode:    http.StatusOK,
	}

	srv := httptest.NewServer(a.Handle(ra))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	req.Header.Set("x-api-key", "valid-key")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status: %d", resp.StatusCode)
	}
}
