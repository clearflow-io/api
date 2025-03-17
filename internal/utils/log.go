package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the incoming request's method, URL, and remote address.
		logrus.Info(fmt.Sprintf("%s %s", r.Method, r.URL.String()))

		// Serve the request.
		next.ServeHTTP(w, r)

		// Log the duration it took to process the request.
		duration := time.Since(start)
		logrus.Info(fmt.Sprintf("Request processed in %s", duration))
	})
}
