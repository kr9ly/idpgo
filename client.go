package idpgo

type Client interface {
	GetClientId() string
	GetClientSecret() string
	IsValidRedirectUri(uri string) bool
}