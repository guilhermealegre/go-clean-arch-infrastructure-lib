package auth

const (
	JWTSecretKey        = "jwtSecret"
	HeaderAuthorization = "Authorization"
	HeaderSetCookie     = "Set-Cookie"
	CookieJwtHeaderBody = "JwtHeaderBody"
	CookieJwtSignature  = "JwtSignature"
	ClaimUsername       = "username"
	ClaimIdMarket       = "idMarket"
	ClaimIdShop         = "idShop"
	ClaimIdBu           = "idBu"
	ClaimIdFascia       = "idFascia"
	ClaimFirstName      = "firstName"
	ClaimLastName       = "lastName"
	ClaimIdLanguage     = "idLanguage"
	ClaimLanguageCode   = "languageCode"
	ClaimIdUserExternal = "idUserExternal"
	ClaimAuthorizations = "authorizations"
)

const (
	TokenTTLMinutes = 30
)
