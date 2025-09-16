package middleware

type authenticationClient interface {
}

type AuthenticationMiddleware struct {
	authClient authenticationClient
}

func NewAuthenticationMiddleware() {}
