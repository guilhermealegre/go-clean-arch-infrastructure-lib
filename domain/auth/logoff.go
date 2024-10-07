package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"bitbucket.org/asadventure/be-core-lib/response"
)

func LogOff(ctx *gin.Context) {
	// validate Signature cookie
	signature, err := ctx.Request.Cookie(CookieJwtSignature)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}

	// will force the client cookie to be deleted
	signature.Expires = time.Unix(0, 0)
	signature.Path = "/"

	// set the cookie to the response
	http.SetCookie(ctx.Writer, signature)

	// validate Header Body Cookie header
	headerBody, err := ctx.Request.Cookie(CookieJwtHeaderBody)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorUnauthorized.Formats(err))
		ctx.Abort()
		return
	}
	// will force the client cookie to be deleted
	headerBody.Expires = time.Unix(0, 0)
	headerBody.Path = "/"

	// set the cookie to the response
	http.SetCookie(ctx.Writer, headerBody)
}
