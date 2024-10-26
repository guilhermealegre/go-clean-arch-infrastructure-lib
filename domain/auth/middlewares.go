package auth

import (
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/errors"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/context"

	"github.com/gin-gonic/gin"
	"github.com/guilhermealegre/go-clean-arch-core-lib/response"
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
	if ctx.Request.URL.Path == "/v1/user/logoff" {
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

	// set token useful claims to context
	internalCtx.SetIdMarket(int(handler.Claims[ClaimIdMarket].(float64)))
	internalCtx.SetIdBu(int(handler.Claims[ClaimIdBu].(float64)))
	internalCtx.SetIdShop(int(handler.Claims[ClaimIdShop].(float64)))
	internalCtx.SetIdFascia(int(handler.Claims[ClaimIdFascia].(float64)))
	internalCtx.SetLanguageCode(handler.Claims[ClaimLanguageCode].(string))
	internalCtx.SetIdUserExternal(int(handler.Claims[ClaimIdUserExternal].(float64)))
	internalCtx.SetUsername(handler.Claims[ClaimUsername].(string))
	authorizations := make([]string, 0)
	if reflect.TypeOf(handler.Claims[ClaimAuthorizations]).Kind() == reflect.Slice &&
		len(handler.Claims[ClaimAuthorizations].([]interface{})) > 0 {
		for _, val := range handler.Claims[ClaimAuthorizations].([]interface{}) {
			authorizations = append(authorizations, val.(string))
		}
	}
	internalCtx.SetAuthorizations(authorizations)

	ctx.Next()
}

func LoadShopParameters(ctx *gin.Context) {
	internalCtx := context.NewContext(ctx)

	if idMarket := ctx.Query(context.ParamIdMarket); idMarket != "" {
		idMarketInt, err := strconv.Atoi(idMarket)
		if err == nil && idMarketInt > 0 {
			internalCtx.SetIdMarket(idMarketInt)
		}
	}

	if idBu := ctx.Query(context.ParamIdBu); idBu != "" {
		idBuInt, err := strconv.Atoi(idBu)
		if err == nil && idBuInt > 0 {
			internalCtx.SetIdBu(idBuInt)
		}
	}

	if idShop := ctx.Query(context.ParamIdShop); idShop != "" {
		idShopInt, err := strconv.Atoi(idShop)
		if err == nil && idShopInt > 0 {
			internalCtx.SetIdShop(idShopInt)
		}
	}

	if idFascia := ctx.Query(context.ParamIdFascia); idFascia != "" {
		idFasciaInt, err := strconv.Atoi(idFascia)
		if err == nil && idFasciaInt > 0 {
			internalCtx.SetIdFascia(idFasciaInt)
		}
	}

	if languageCode := ctx.Query(context.ParamLanguageCode); languageCode != "" {
		if len(languageCode) > 0 {
			internalCtx.SetLanguageCode(languageCode)
		}

	}

	ctx.Next()
}

func LoadJWTSecret(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(JWTSecretKey, secret)

		ctx.Next()
	}
}
