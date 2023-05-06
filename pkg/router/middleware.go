package router

import "net/http"

// HealthMiddleware is a middleware that responds to health checks
func HealthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SignatureMiddleware is a middleware that checks the signature of the request against the request body
func SignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// GET requests are not signed
			next.ServeHTTP(w, r)
			return
		}

		// TODO implement signature check

		next.ServeHTTP(w, r)
	})
}
