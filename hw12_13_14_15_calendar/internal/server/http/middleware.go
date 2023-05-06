package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server"
)

func loggingMiddleware(logger server.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Since(timeStart)

		ctx := r.Context()
		statusCode, _ := ctx.Value(statusCodeKey).(int)

		remoteAddrParts := strings.SplitN(r.RemoteAddr, ":", 2)

		// 66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
		logger.Infof(
			"%s [%s %s] %s %s HTTP/%d.%d %d %dms %q\n",
			remoteAddrParts[0],
			timeStart.Format("01/Jan/2006"),
			timeStart.Format("15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.ProtoMajor,
			r.ProtoMinor,
			statusCode,
			elapsed.Milliseconds(),
			r.UserAgent(),
		)
	})
}
