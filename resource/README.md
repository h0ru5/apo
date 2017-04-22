# Example resource server

This is an example for a resource server (RS) 
using A.P.O. as identity and access management (IAM) server

It exposes two protected resources called `/open` and `/close` and an unprotected index under the root resource
 
when its using the default configuration, the example RS will try to find A.P.O. on `localhost:3000`
you can  provide commandline params, a configuration file or an environment variable to change that.