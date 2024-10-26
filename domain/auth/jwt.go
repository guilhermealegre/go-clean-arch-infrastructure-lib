package auth

import (
	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain/context"
	"net/http"
	"strings"
	"time"

	"github.com/guilhermealegre/go-clean-arch-core-lib/helpers"
	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/errors"
)

type JWTHandler struct {
	SecretKey   string        `json:"secretKey"`
	ExpiryTime  int           `json:"expiryTime"`
	Claims      jwt.MapClaims `json:"claims"`
	TokenString string        `json:"tokenString"`
}

func NewJWTHandler(secretKey string, expiryTime int) *JWTHandler {
	return &JWTHandler{
		SecretKey:  secretKey,
		ExpiryTime: expiryTime,
		Claims:     jwt.MapClaims{},
	}
}

func NewAuthorizationHandler(authorizationHeader string, secretKey string) (*JWTHandler, error) {
	// Verify if its empty
	if authorizationHeader == "" {
		return nil, errors.ErrorAuthorizationMissing()
	}
	// Split the string to get rid of the Bearer part
	// and get the jwt string
	var tokenString string
	if aToken := strings.Split(authorizationHeader, " "); len(aToken) > 1 {
		if helpers.TrimAndLowerStr(aToken[0]) != "bearer" {
			return nil, errors.ErrorInvalidBearerKey()
		}
		tokenString = aToken[1]
	} else {
		return nil, errors.ErrorInvalidBearerKey()
	}
	// Parse token string to get Claims
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrorUnexpectedSigninMethod().Formats(token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return &JWTHandler{
		SecretKey:   secretKey,
		TokenString: tokenString,
		Claims:      *token.Claims.(*jwt.MapClaims),
	}, nil
}

func (j *JWTHandler) GenerateToken() (err error) {

	// Verify if token has claims
	if len(j.Claims) == 0 {
		return errors.ErrorNoClaimsFound()
	}

	// Add token expiry time
	j.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(j.ExpiryTime)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.Claims)
	j.TokenString, err = token.SignedString([]byte(j.SecretKey))
	return err
}

func (j *JWTHandler) AddClaim(key string, value any) *JWTHandler {
	j.Claims[key] = value
	return j
}

func (j *JWTHandler) SetSplitCookies(gCtx contextDomain.IContext, maxAge int) error {

	expirationTime := time.Now().Add(TokenTTLMinutes * time.Minute)
	// validate Token string
	if j.TokenString == "" {
		return errors.ErrorTokenStringEmpty()
	}

	// Create a cookie with 30-minute expiration
	cookieHeaderBody := &http.Cookie{
		Name:     CookieJwtHeaderBody,
		Value:    j.HeaderBodyPart(),
		Expires:  expirationTime,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	// Set the cookie to the response
	http.SetCookie(gCtx.Response(), cookieHeaderBody)

	// Create a cookie with 30-minute expiration
	cookieSignature := &http.Cookie{
		Name:     CookieJwtSignature,
		Value:    j.SignaturePart(),
		Expires:  expirationTime,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	// Set the cookie to the response
	http.SetCookie(gCtx.Response(), cookieSignature)

	return nil
}

func (j *JWTHandler) HeaderBodyPart() string {
	return j.TokenString[:strings.LastIndex(j.TokenString, ".")]
}

func (j *JWTHandler) SignaturePart() string {
	return j.TokenString[strings.LastIndex(j.TokenString, ".")+1:]
}
