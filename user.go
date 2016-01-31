package idpgo

import "net/http"

type User interface {
	GetUserId() string
	CredentialMatch(r *http.Request) bool
}
