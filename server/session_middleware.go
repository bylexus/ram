package server

import (
	"context"
	"net/http"

	"github.com/bylexus/go-stdlib/log"
	"github.com/kataras/go-sessions/v3"
)

type SessionMiddleware interface {
	WrapHandler(handler http.Handler) http.Handler
}

type concreteSessionMiddleware struct {
	logger         *log.SeverityLogger
	sessionHandler sessions.Sessions
}

const SessionContextKey ServerContextKey = "session"

type handlerWrapper struct {
	middlewareInst *concreteSessionMiddleware
	wrappedHandler http.Handler
}

func (h handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.middlewareInst.serveHttpFn(h.wrappedHandler, w, r)
}

func NewSessionMiddleware(logger *log.SeverityLogger) SessionMiddleware {
	s := sessions.New(sessions.Config{})
	return concreteSessionMiddleware{
		logger:         logger,
		sessionHandler: *s,
	}
}

func (s concreteSessionMiddleware) WrapHandler(handler http.Handler) http.Handler {
	h := handlerWrapper{
		middlewareInst: &s,
		wrappedHandler: handler,
	}
	return h
}

func (s concreteSessionMiddleware) serveHttpFn(h http.Handler, w http.ResponseWriter, r *http.Request) {

	// Start session and inject it to the request context
	// get in child requests with r.Context().Value(SessionContextKey)
	session := s.sessionHandler.Start(w, r)
	ctx := context.WithValue(r.Context(), SessionContextKey, session)

	// call original handler:
	h.ServeHTTP(w, r.WithContext(ctx))
}

func GetSession(r *http.Request) *sessions.Session {
	return r.Context().Value(SessionContextKey).(*sessions.Session)
}
