package http

import (
	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/domain/auth"
	"net/url"

	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/domain"

	"github.com/gin-gonic/gin"
)

// Endpoint Struct
type Endpoint struct {
	// Group
	group *Group
	// Path
	path string
	// Method
	method string
	// Middleware
	middlewares []gin.HandlerFunc
	// authorizations
	authorizations []string
}

// NewEndpoint creates a new endpoint
func NewEndpoint(group *Group, path string, method string) *Endpoint {
	return &Endpoint{
		group:       group,
		path:        path,
		method:      method,
		middlewares: make([]gin.HandlerFunc, 0),
	}
}

// Path endpoint path
func (e *Endpoint) Path() string {
	return e.path
}

// Method http method
func (e *Endpoint) Method() string {
	return e.method
}

// AddMiddlewares sets middlewares with handle functions
func (e *Endpoint) AddMiddlewares(middlewares domain.IMiddleware) {
	e.middlewares = append(e.middlewares, middlewares.GetHandlers()...)
}

// RequireAuthorizations sets the route required authorizations
func (e *Endpoint) RequireAuthorizations(authorizations ...string) *Endpoint {
	e.authorizations = append(e.authorizations, authorizations...)
	return e
}

// SetRoute sets a route with handle functions
func (e *Endpoint) SetRoute(engine *gin.Engine, handlerFunc ...gin.HandlerFunc) {
	if len(e.authorizations) > 0 {
		e.middlewares = append(e.middlewares, auth.NewAuthorizationCheck(e.authorizations...))
	}
	e.middlewares = append(e.middlewares, handlerFunc...)
	e.group.Init(&engine.RouterGroup).Handle(e.Method(), e.Path(), e.middlewares...)
}

// FullPath full endpoint path
func (e *Endpoint) FullPath() string {
	path, _ := url.JoinPath(e.group.String(), e.path)
	return path
}
