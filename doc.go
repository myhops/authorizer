/*
Authorizer is a very simple web server that checks api keys passed 
in http headers or a request.
Its main use case is to use it as an Istio HTTP Authorizer 
and enable the save use of API keys in Authorization Policies.

*/
package authorizer