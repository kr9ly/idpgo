package idpgo

import "net/http"

type Service interface {
	GetClient(id string) (Client, error)
	GetUser(r *http.Request) (User, error)
	NewAccessToken(params *AuthorizationParams, client Client, user User) (string, error)
	NewAuthorizationCode(accessToken string, idToken string) (string, error)
}
