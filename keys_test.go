package authorizer

import (
	"net/http"
	"testing"
)


func TestKeys_authorizeHeader(t *testing.T) {
	type args struct {
		h http.Header
	}
	tests := []struct {
		name    string
		k       Keys
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			k: Keys{
				{
					Header: http.CanonicalHeaderKey("key-header"),
					Values: []string{"value1"},
				},
			},
			args: args {
				h: http.Header(map[string][]string{
					"Key-Header": {"value1"},
				}),
			},
			wantErr: false,
		},
		{
			name: "err",
			k: Keys{
				{
					Header: http.CanonicalHeaderKey("key-header"),
					Values: []string{"value1"},
				},
			},
			args: args {
				h: http.Header(map[string][]string{
					"Key-Header": {"value2"},
				}),
			},
			wantErr: true,
		},

		// 
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.k.authorizeHeader(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("Keys.authorizeHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
