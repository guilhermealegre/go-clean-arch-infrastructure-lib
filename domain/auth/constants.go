package auth

const (
	JWTSecretKey        = "jwtSecret"
	HeaderAuthorization = "Authorization"
	HeaderSetCookie     = "Set-Cookie"
	CookieJwtHeaderBody = "JwtHeaderBody"
	CookieJwtSignature  = "JwtSignature"
	ClaimIdUser         = "id_user"
	ClaimFirstName      = "first_name"
	ClaimLastName       = "last_name"
	ClaimEmail          = "email"
	ClaimPhoneNumber    = "phone_number"
	ClaimAuthorizations = "authorizations"
)

const (
	TokenTTLMinutes = 30
)
