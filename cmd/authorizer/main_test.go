package main

import (
	"log/slog"
	"reflect"
	"testing"
)

func Test_getOptions(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *options
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				[]string{"app", "-key=header1=key1", "-key=header2=key2"},
			},
			want: &options{
				ListenAddress: ":8080",
				Keys: keysOption{
					keyOption{
						Header: "header1",
						Key: "key1",
					},
					keyOption{
						Header: "header2",
						Key: "key2",
					},
				},
				AllowedCode: 200,
				ForbiddenCode: 403,
				LogFile: "/dev/stdout",
				LogFormat: "text",
			},
		},
		{
			name: "level",
			args: args{
				[]string{"app", "-key=header1=key1", "-key=header2=key2", "-loglevel=warn"},
			},
			want: &options{
				ListenAddress: ":8080",
				Keys: keysOption{
					keyOption{
						Header: "header1",
						Key: "key1",
					},
					keyOption{
						Header: "header2",
						Key: "key2",
					},
				},
				AllowedCode: 200,
				ForbiddenCode: 403,
				LogLevel: logLevelOption(slog.LevelWarn),
				LogFile: "/dev/stdout",
				LogFormat: "text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOptions(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
