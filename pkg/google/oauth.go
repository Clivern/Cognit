// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	OauthAuthorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	OauthTokenURL     = "https://oauth2.googleapis.com/token"
	OauthUserURL      = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// OAuthConfig holds Google OAuth app credentials.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuth exchanges authorization codes for access tokens.
type OAuth struct {
	cfg OAuthConfig
}

// Token is an OAuth access token response.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IDToken     string `json:"id_token"`
}

// UserInfo is the authenticated Google user.
type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

// NewOAuth returns an OAuth helper for the given app configuration.
func NewOAuth(cfg OAuthConfig) *OAuth {
	return &OAuth{cfg: cfg}
}

// AuthorizeURL builds the Google authorization URL for the provided state.
func (o *OAuth) AuthorizeURL(state string) string {
	q := url.Values{
		"client_id":     {o.cfg.ClientID},
		"redirect_uri":  {o.cfg.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(o.cfg.Scopes, " ")},
		"state":         {state},
		"access_type":   {"online"},
		"prompt":        {"select_account"},
	}

	return OauthAuthorizeURL + "?" + q.Encode()
}
