package auth

import (
	"encoding/json"
	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/domain"
	"net/http"
	"time"

	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/errors"

	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/context"

	"bitbucket.org/asadventure/be-core-lib/response"
	"github.com/gin-gonic/gin"
)

func BuildAuthorizationHeader(ctx *gin.Context) {
	// validate Signature cookie
	signature, err := ctx.Cookie(CookieJwtSignature)
	if err != nil {
		// return unauthorized error
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}
	if signature == "" {
		// return unauthorized error
		ctx.JSON(http.StatusUnauthorized, errors.ErrorMissingSignatureCookie())
		ctx.Abort()
		return
	}

	// validate authorization header
	authorizationString := ctx.GetHeader(HeaderAuthorization)

	// validate emptiness
	if authorizationString != "" && signature != "" {
		authorizationString += "." + signature
		// complete Authorization
		ctx.Request.Header.Set(HeaderAuthorization, authorizationString)
		ctx.Next()
		return
	} else {
		// return unauthorized error
		ctx.JSON(http.StatusUnauthorized, errors.ErrorAuthorizationMissing())
		ctx.Abort()
		return
	}
}

func IncreaseActivityTTLInXMinutes(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/v1/auth/logoff" {
		return
	}

	// Read the cookie from the request
	cookieSignature, err := ctx.Cookie(CookieJwtSignature)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	// Create a new cookie with an extended TTL
	newSignatureCookie := &http.Cookie{
		Name:     CookieJwtSignature,
		Value:    cookieSignature,
		Expires:  time.Now().Add(TokenTTLMinutes * time.Minute),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	// Set the cookie to the response
	http.SetCookie(ctx.Writer, newSignatureCookie)

	// Read the cookie from the request
	cookieHeaderBody, err := ctx.Cookie(CookieJwtHeaderBody)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	// Create a new cookie with an extended TTL
	newHeaderBodyCookie := &http.Cookie{
		Name:     CookieJwtHeaderBody,
		Value:    cookieHeaderBody,
		Expires:  time.Now().Add(TokenTTLMinutes * time.Minute),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	// Set the cookie to the response
	http.SetCookie(ctx.Writer, newHeaderBodyCookie)

	ctx.Next()
}

func ValidateToken(ctx *gin.Context) {
	internalCtx := context.NewContext(ctx)
	jwtAuthorizationHeader := ctx.GetHeader(HeaderAuthorization)
	secret := internalCtx.GetString(JWTSecretKey)

	if secret == "" {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats("JWT Secret is not defined"))
		ctx.Abort()
		return
	}

	handler, err := NewAuthorizationHandler(jwtAuthorizationHeader, secret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	user := &domain.UserDetailJwt{}
	claimsBytes, err := json.Marshal(handler.Claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	err = json.Unmarshal(claimsBytes, user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	internalCtx.SetUser(user)

	ctx.Next()
}

func LoadJWTSecret(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(JWTSecretKey, secret)

		ctx.Next()
	}
}
