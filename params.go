package idpgo

import "net/url"

type AuthorizationParams struct {
	HasCode bool
	HasToken bool
	HasIdToken bool
	RedirectUri *url.URL
	Scopes []string
	State string
	Nonce string
}
