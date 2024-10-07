package auth

import (
	"bitbucket.org/asadventure/be-core-lib/response"
	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/context"
	"net/http"
)

type AuthorizationCheck struct {
	authorizations []string
}

func NewAuthorizationCheck(authorizations ...string) gin.HandlerFunc {
	middleware := &AuthorizationCheck{
		authorizations: authorizations,
	}

	return middleware.check
}

func (m *AuthorizationCheck) check(ctx *gin.Context) {
	internalCtx := context.NewContext(ctx)
	enc := 0
	for _, required := range m.authorizations {
		for _, authorization := range internalCtx.GetAuthorizations() {
			if authorization == required {
				enc++
				break
			}
		}
	}

	if enc != len(m.authorizations) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats("You don't have access to this resource"))
		ctx.Abort()
		return
	}

	ctx.Next()
}
