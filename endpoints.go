package idpgo

import (
	"net/http"
	"html/template"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthorizationEndpoint struct {
	Issuer        string
	Service       Service
	Template      *template.Template
	SigningMethod jwt.SigningMethod
	SigningKey    interface{}
	StrictMode    bool
}

func (e *AuthorizationEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &RequestWrapper{r}
	resp := &ResponseWrapper{w}
	switch r.Method {
	case "GET":
		e.displayAuthorizationView(resp, req)
		return
	case "POST":
		e.authorizeUser(resp, req)
		return
	}
	resp.MethodNotAllowed("method not allowed.")
}

func (e *AuthorizationEndpoint) displayAuthorizationView(resp *ResponseWrapper, req *RequestWrapper) {
	if _, _, err := req.StartAuthorization(e.Service, e.StrictMode); err != nil {
		resp.BadRequest(err.Error())
		return
	}

	resp.ResponseTemplate(e.Template, req)
}

func (e *AuthorizationEndpoint) authorizeUser(resp *ResponseWrapper, req *RequestWrapper) {
	var params *AuthorizationParams
	var client Client
	var err error
	if params, client, err = req.StartAuthorization(e.Service, e.StrictMode); err != nil {
		resp.BadRequest(err.Error())
		return
	}
	var user User
	if user, err = e.Service.GetUser(req.Request); err != nil {
		resp.ResponseTemplate(e.Template, req)
		return
	}

	if !user.CredentialMatch(req.Request) {
		resp.ResponseTemplate(e.Template, req)
		return
	}

	redirectTo := params.RedirectUri
	var accessToken string
	if accessToken, err = e.Service.NewAccessToken(params, client, user); err != nil {
		resp.ResponseTemplate(e.Template, req)
		return
	}

	var idToken string
	if params.HasIdToken {
		token := jwt.New(e.SigningMethod)
		token.Claims["iss"] = e.Issuer
		token.Claims["user_id"] = user.GetUserId()
		token.Claims["aud"] = client.GetClientId()
		token.Claims["iat"] = time.Now().Add(600 * time.Second).Unix()
		token.Claims["exp"] = time.Now()
		token.Claims["nonce"] = params.Nonce

		var signed string
		if signed, err = token.SignedString(e.SigningKey); err != nil {
			resp.BadRequest("jwt error:" + err.Error())
			return
		}
		idToken = signed
	}

	if params.HasCode {
		q := Query{}
		if params.State != "" {
			q["state"] = params.State
		}

		var authorizationCode string
		if authorizationCode, err = e.Service.NewAuthorizationCode(accessToken, idToken); err != nil {
			resp.ResponseTemplate(e.Template, req)
			return
		}

		q["code"] = authorizationCode
		redirectTo.RawQuery = q.String()
	}

	if params.HasToken {
		q := Query{}
		q["access_token"] = accessToken
		q["id_token"] = idToken
		q["token_type"] = "bearer"
		if params.State != "" {
			q["state"] = params.State
		}
		q["expires_in"] = "3600"

		redirectTo.Fragment = q.String()
	}

	resp.SeeOther(redirectTo.String())
}