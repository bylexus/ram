package server

import (
	"errors"
	"net/http"

	"github.com/bylexus/go-stdlib/log"
)

/*
TODO: Document AuthMiddleware
*/
type AuthMiddleware interface {
	WrapHandler(handler http.Handler) http.Handler
}

type concreteAuthMiddleware struct {
	logger          *log.SeverityLogger
	exceptionRoutes []string
}

const userContextKey = "user"

type authHandlerWrapper struct {
	middlewareInst *concreteAuthMiddleware
	wrappedHandler http.Handler
}

func (h authHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.middlewareInst.serveHttpFn(h.wrappedHandler, w, r)
}

func NewAuthMiddleware(logger *log.SeverityLogger) AuthMiddleware {
	return concreteAuthMiddleware{
		logger:          logger,
		exceptionRoutes: make([]string, 0),
	}
}

func (s concreteAuthMiddleware) WrapHandler(handler http.Handler) http.Handler {
	h := authHandlerWrapper{
		middlewareInst: &s,
		wrappedHandler: handler,
	}
	return h
}

func (s concreteAuthMiddleware) serveHttpFn(h http.Handler, w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("auth for route: %s", r.RequestURI)

	session := GetSession(r)
	if session == nil {
		NewErrorJsonResponse(nil, http.StatusInternalServerError, errors.New("no session"), http.StatusInternalServerError).WriteHttpResponse(w)
		return
	}
	user := r.Context().Value(userContextKey)
	if user == nil {
		http.Redirect(w, r, "/guest/login.html", http.StatusTemporaryRedirect)
		return
	}

	// call original handler:
	h.ServeHTTP(w, r)
}
