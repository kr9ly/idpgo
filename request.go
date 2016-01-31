package idpgo

import (
	"net/http"
	"strings"
	"net/url"
)

type RequestWrapper struct {
	*http.Request
}

func (req *RequestWrapper) StartAuthorization(db Service, strictMode bool) (*AuthorizationParams, Client, error) {
	rType := req.FormValue("response_type")
	cId := req.FormValue("client_id")
	uri := req.FormValue("redirect_uri")
	scope := req.FormValue("scope")
	state := req.FormValue("state")
	nonce := req.FormValue("nonce")

	if rType == "" {
		return nil, nil, BadRequestError("response_type must be required.")
	}

	if cId == "" {
		return nil, nil, BadRequestError("client_id must be required.")
	}

	if uri == "" {
		return nil, nil, BadRequestError("redirect_uri must be required.")
	}

	if uri == "" {
		return nil, nil, BadRequestError("redirect_uri must be required.")
	}

	if scope == "" {
		return nil, nil, BadRequestError("scope must be required.")
	}

	hasCode := false
	hasToken := false
	hasIdToken := false

	responseTypes := strings.Split(rType, " ")
	for i := range responseTypes {
		switch responseTypes[i] {
		case "code":
			hasCode = true
		case "token":
			hasToken = true
		case "id_token":
			hasIdToken = true
		}
	}

	hasOpenIdScope := false
	scopes := strings.Split(scope, " ")
	for i := range scopes {
		if scopes[i] == "openid" {
			hasOpenIdScope = true
		}
	}

	if !hasCode && !hasToken && !hasIdToken {
		return nil, nil, BadRequestError("response_type must be contains 'code' or 'token' or 'id_token'.")
	}

	if hasIdToken && !hasOpenIdScope {
		return nil, nil, BadRequestError("openid scope is required on openid connect request.")
	}

	if strictMode {
		if state == "" {
			return nil, nil, BadRequestError("state must be required.")
		}

		if hasIdToken && nonce == "" {
			return nil, nil, BadRequestError("nonce must be required.")
		}
	}

	var client Client
	var err error
	if client, err = db.GetClient(cId); err != nil {
		return nil, nil, BadRequestError("client_id is invalid.")
	}

	if !client.IsValidRedirectUri(uri) {
		return nil, nil, BadRequestError("redirect_uri is invalid.")
	}

	var parsedUri *url.URL
	if parsedUri, err = url.Parse(uri); err != nil {
		return nil, nil, BadRequestError("redirect_uri is invalid.")
	}

	return &AuthorizationParams{
		HasCode: hasCode,
		HasToken: hasToken,
		HasIdToken: hasIdToken,
		RedirectUri: parsedUri,
		Scopes: scopes,
		State: state,
		Nonce: nonce,
	}, client, nil
}