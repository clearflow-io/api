package api

import (
	"fmt"
	"net/http"
)

// EnforceHTTPS redirects HTTP requests to HTTPS if the environment is not "local".
func EnforceHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the X-Forwarded-Proto header set by proxies (Render, Railway, AWS, etc.)
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "http" {
			target := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
