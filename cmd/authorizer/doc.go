/*
Authorizer is a simple web server that checks the header of the
http request for the existence of valid api keys.

	$ ./authorizer -h
	Usage of ./authorizer:
	  -address string
		Listen address, defaults to :8080 (default ":8080")
	  -allowed-code int
		status code for allowed access (default 200)
	  -forbidden-code int
		status code for forbidden access (default 403)
	  -key value
		valid headers, <header,value>. Can be used multiple times.
	  -logformat value
		logformat: [json, text] (default text)
	  -loglevel value
		loglevel: [debug, info, warn, error]

The keys can be specified using multiple flags or you can specify more than one key per flag.
	$ ./authorizer \
		-key=header=value,header=value1 \
		-key=header=value3 \
		-key=header2=value2


*/
package main
