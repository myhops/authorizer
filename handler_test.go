package authorizer

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizer_Authorized(t *testing.T) {
	type fields struct {
		Keys                    []*Key
		NotAuthorizedStatusCode int
	}
	type args struct {
		h http.Header
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Keys: []*Key{
					{
						Header: http.CanonicalHeaderKey("key-header"),
						Values: []string{"value1"},
					},
				},
				NotAuthorizedStatusCode: http.StatusForbidden,
			},
			args: args{
				h: http.Header(map[string][]string{
					"Key-Header": {"value1"},
				}),
			},
			want: true,
		},
		{
			name: "not ok",
			fields: fields{
				Keys: []*Key{
					{
						Header: http.CanonicalHeaderKey("key-header"),
						Values: []string{"value1"},
					},
				},
				NotAuthorizedStatusCode: http.StatusForbidden,
			},
			args: args{
				h: http.Header(map[string][]string{
					"Key-Header": {"value2"},
				}),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Authorizer{
				Keys:                    tt.fields.Keys,
				NotAuthorizedStatusCode: tt.fields.NotAuthorizedStatusCode,
			}
			if got := a.Authorized(tt.args.h); got != tt.want {
				t.Errorf("Authorizer.Authorized() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestHandler(t *testing.T) {
	newKey := func(key string, values []string) *Key {
		return &Key{
			Header: http.CanonicalHeaderKey(key),
			Values: values,
		}
	}

	a := &Authorizer{
		Keys: []*Key{
			newKey("x-api-key", []string{"valid-key"}),
		},
		Logger: slog.Default(),
		NotAuthorizedStatusCode: http.StatusForbidden,
		AuthorizedStatusCode: http.StatusOK,
	}

	srv := httptest.NewServer(a.Handle())
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