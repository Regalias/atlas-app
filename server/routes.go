package server

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func (s *server) routes(appLogger *zerolog.Logger) {
	// Setup middleware chain
	// Build middleware chains from base logger
	c := alice.New().Append(hlog.NewHandler(*appLogger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	// c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	s.router.Handler("GET", "/", c.ThenFunc(s.handleIndex()))
	s.router.Handler("GET", "/:linkpath", c.ThenFunc(s.handleServeLink()))

}
