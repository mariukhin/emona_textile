package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		if token == nil || !token.Valid {
			httpResponseError(w, http.StatusUnauthorized, "request__unauthorized")
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func AllowContentTypeJson() func(next http.Handler) http.Handler {
	return AllowContentType("application/json")
}

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
func AllowContentType(contentTypes ...string) func(next http.Handler) http.Handler {
	cT := []string{}
	for _, t := range contentTypes {
		cT = append(cT, strings.ToLower(t))
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			for _, t := range cT {
				if t == s {
					next.ServeHTTP(w, r)
					return
				}
			}

			httpResponseError(w, http.StatusUnsupportedMediaType, "content_type__invalid")
		}
		return http.HandlerFunc(fn)
	}
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				//logEntry := GetLogEntry(r)
				//if logEntry != nil {
				//	logEntry.Panic(rvr, debug.Stack())
				//} else {
				//	middleware.PrintPrettyStack(rvr)
				//}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func RequestLogger(appLogger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				remoteIP = r.RemoteAddr
			}

			requestLogger := appLogger.WithFields(logrus.Fields{
				"request_id": requestId,
				"remote_ip":  remoteIP,
			})

			next.ServeHTTP(w, r.WithContext(logger.WithLogger(r.Context(), requestLogger)))
		}
		return http.HandlerFunc(fn)
	}
}

func HttpLogger(logger *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := logrus.Fields{
					"status":    ww.Status(),
					"bytes":     ww.BytesWritten(),
					"duration":  int64(time.Since(t1)),
					"durationf": time.Since(t1).String(),
					"remote_ip": remoteIP,
					"proto":     r.Proto,
					"method":    r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				logger.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
