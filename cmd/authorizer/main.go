package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

type keyOption struct {
	Header string
	Key    string
}

type options struct {
	ListenAddress string
	Keys          keysOption
	LogLevel      logLevelOption
	LogFormat     logFormatOption
	LogFile       string
	AllowedCode   int
	ForbiddenCode int
}

type logFormatOption string
 
func (l *logFormatOption) Set(value string) error {
	switch strings.ToLower(value) {
	case "text", "json":
		*l = logFormatOption( value)
	default:
		return fmt.Errorf("bad log format: %s", value)
	}
	return nil
}

func (l *logFormatOption) String() string {
	return string(*l)
}

type logLevelOption slog.Level

func (l *logLevelOption) Set(value string) error {
	switch strings.ToUpper(value) {
	case "ERROR":
		*l = logLevelOption(slog.LevelError)
	case "INFO":
		*l = logLevelOption(slog.LevelInfo)
	case "WARN":
		*l = logLevelOption(slog.LevelWarn)
	case "DEBUG":
		*l = logLevelOption(slog.LevelDebug)
	default:
		return fmt.Errorf("bad value for log level: %s", value)
	}
	return nil
}

func (l *logLevelOption) String() string {
	return slog.Level(*l).String()
}

type keysOption []keyOption

func (o *keysOption) Set(value string) error {
	// Split the value in key and header.
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return fmt.Errorf("key needs to have a key and a value separated by a comma")
	}
	// Append the key
	*o = append(*o, keyOption{Header: s[0], Key: s[1]})
	return nil
}

func (o *keysOption) String() string {
	return fmt.Sprint(*o)
}

func getOptions(args []string) (*options, error) {
	opts := &options{}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	args = args[1:]

	fs.Var(&opts.Keys, "key", "valid headers, <header,value>. Can be used multiple times.")
	fs.StringVar(&opts.ListenAddress, "address", ":8080", "Listen address, defaults to :8080")
	fs.Var(&opts.LogLevel, "loglevel", "loglevel: [debug, info, warn, error]")
	fs.Var(&opts.LogFormat, "logformat", "logformat: [json, text]")
	fs.IntVar(&opts.AllowedCode, "allowed-code", 200, "status code for allowed access")
	fs.IntVar(&opts.ForbiddenCode, "forbidden-code", 403, "status code for forbidden access")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opts, nil
}

func main() {
	options, err := getOptions(os.Args)
	if err != nil {
		log.Printf("error: %s", err)
	}
	_ = options
}
