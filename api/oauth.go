// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/cognit/db"
	"github.com/clivern/cognit/module"
	"github.com/clivern/cognit/pkg/github"
	"github.com/clivern/cognit/pkg/google"
	"github.com/clivern/cognit/pkg/resend"
	"github.com/clivern/cognit/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const OauthStateCookie = "_cognit_oauth_state"

// GitHubOAuthStartAction redirects the browser to GitHub's authorize URL.
func GitHubOAuthStartAction(w http.ResponseWriter, r *http.Request) {
	oauth := github.NewOAuth(github.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.github.client_id"),
		ClientSecret: viper.GetString("app.oauth.github.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.github.redirect_url"),
		Scopes:       []string{"read:user", "user:email"},
		AllowSignup:  true,
	})

	state, err := util.GenerateSecureToken(24)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate oauth state")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=github"), http.StatusFound)
		return
	}

	authorizeURL := oauth.AuthorizeURL(state)

	opts := lo.Ternary(
		strings.HasPrefix(util.AppURL(""), "https://"),
		util.SecureCookieOptions(),
		util.DefaultCookieOptions(),
	)
	opts.SameSite = http.SameSiteLaxMode
	opts.MaxAge = int((10 * time.Minute) / time.Second)
	util.SetCookie(w, OauthStateCookie, state, opts)

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

// GitHubOAuthCallbackAction completes GitHub OAuth and creates a session.
func GitHubOAuthCallbackAction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	expectedState := util.GetCookie(r, OauthStateCookie)

	util.DeleteCookie(w, OauthStateCookie)

	oauth := github.NewOAuth(github.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.github.client_id"),
		ClientSecret: viper.GetString("app.oauth.github.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.github.redirect_url"),
		Scopes:       []string{"read:user", "user:email"},
		AllowSignup:  true,
	})

	token, err := oauth.Exchange(r.Context(), code, state, expectedState)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth exchange failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=github"), http.StatusFound)
		return
	}

	user, err := oauth.User(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth user fetch failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=github"), http.StatusFound)
		return
	}

	emails, err := oauth.Emails(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth emails fetch failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=github"), http.StatusFound)
		return
	}

	auth := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
	)

	result, err := auth.LoginWithOAuth(r.Context(), &module.OAuthIdentity{
		Provider:       db.UserProviderGithub,
		ProviderUserID: strconv.FormatInt(user.ID, 10),
		Email:          github.PrimaryEmail(emails, user.Email),
		Name: lo.Ternary(
			lo.IsNotEmpty(user.Name),
			user.Name,
			user.Login,
		),
	})
	if err != nil {
		log.Error().Err(err).Msg("GitHub oauth login failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=github"), http.StatusFound)
		return
	}

	util.SetCookie(w, "_cognit_session", result.Session.Token, result.CookieOptions)

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	err = im.AttachPending(result.User)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", result.User.Id.String()).
			Msg("Failed to attach pending workspace invites")
	}

	http.Redirect(w, r, util.AppURL("/login?oauth=github"), http.StatusFound)
}

// GoogleOAuthStartAction redirects the browser to Google's authorize URL.
func GoogleOAuthStartAction(w http.ResponseWriter, r *http.Request) {
	oauth := google.NewOAuth(google.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.google.client_id"),
		ClientSecret: viper.GetString("app.oauth.google.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.google.redirect_url"),
		Scopes:       []string{"openid", "email", "profile"},
	})

	state, err := util.GenerateSecureToken(24)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate oauth state")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=google"), http.StatusFound)
		return
	}

	authorizeURL := oauth.AuthorizeURL(state)

	opts := lo.Ternary(
		strings.HasPrefix(util.AppURL(""), "https://"),
		util.SecureCookieOptions(),
		util.DefaultCookieOptions(),
	)
	opts.SameSite = http.SameSiteLaxMode
	opts.MaxAge = int((10 * time.Minute) / time.Second)
	util.SetCookie(w, OauthStateCookie, state, opts)

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

// GoogleOAuthCallbackAction completes Google OAuth and creates a session.
func GoogleOAuthCallbackAction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	expectedState := util.GetCookie(r, OauthStateCookie)

	util.DeleteCookie(w, OauthStateCookie)

	oauth := google.NewOAuth(google.OAuthConfig{
		ClientID:     viper.GetString("app.oauth.google.client_id"),
		ClientSecret: viper.GetString("app.oauth.google.client_secret"),
		RedirectURL:  viper.GetString("app.oauth.google.redirect_url"),
		Scopes:       []string{"openid", "email", "profile"},
	})

	token, err := oauth.Exchange(r.Context(), code, state, expectedState)
	if err != nil {
		log.Error().Err(err).Msg("Google oauth exchange failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=google"), http.StatusFound)
		return
	}

	user, err := oauth.User(r.Context(), token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("Google oauth user fetch failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=google"), http.StatusFound)
		return
	}

	auth := module.NewAuth(
		db.NewUserRepository(db.GetDB()),
		db.NewSessionRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
	)

	result, err := auth.LoginWithOAuth(r.Context(), &module.OAuthIdentity{
		Provider:       db.UserProviderGoogle,
		ProviderUserID: user.Sub,
		Email:          user.Email,
		Name: lo.Ternary(
			lo.IsNotEmpty(user.Name),
			user.Name,
			user.Email,
		),
	})
	if err != nil {
		log.Error().Err(err).Msg("Google oauth login failed")
		http.Redirect(w, r, util.AppURL("/login?oauth_error=google"), http.StatusFound)
		return
	}

	util.SetCookie(w, "_cognit_session", result.Session.Token, result.CookieOptions)

	im := module.NewInvite(
		db.NewUserInviteRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewConfigRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
		db.NewWorkspaceUserRepository(db.GetDB()),
		resend.NewMailer(),
	)
	err = im.AttachPending(result.User)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", result.User.Id.String()).
			Msg("Failed to attach pending workspace invites")
	}

	http.Redirect(w, r, util.AppURL("/login?oauth=google"), http.StatusFound)
}
