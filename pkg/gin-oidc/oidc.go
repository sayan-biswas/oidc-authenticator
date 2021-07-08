package oidc

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/pkg/go-oidc"
	"golang.org/x/oauth2"
)

const contextKey = "github.com/sayan-biswas/gin-oidc"

type Authenticator struct {
	config   *Configuration
	provider *oidc.Provider
	oauth2   *oauth2.Config
	context  *gin.Context
}

type Configuration struct {
	Provider     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func Get(context *gin.Context) *Authenticator {
	return context.MustGet(contextKey).(*Authenticator)
}

func (auth *Authenticator) GetConfiguration() *Configuration {
	return auth.config
}

func (auth *Authenticator) SetContext(context *gin.Context) {
	auth.context = context
	context.Set(contextKey, auth)
}

func New(config *Configuration) (*Authenticator, error) {
	provider, providerError := oidc.NewProvider(oauth2.NoContext, config.Provider)
	if providerError != nil {
		return nil, fmt.Errorf("OIDC Provider Error - %w", providerError)
	}

	oauth2 := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Oauth2Endpoint(),
		Scopes:       config.Scopes,
	}

	authenticator := &Authenticator{
		config,
		provider,
		&oauth2,
		&gin.Context{},
	}

	return authenticator, nil
}

func (auth *Authenticator) LoginHandler() {

	redirectURL, parseError := url.Parse(auth.context.Request.Referer())
	if parseError != nil {
		redirectURL, parseError = url.Parse(auth.context.Request.Host)
		if parseError != nil {
			redirectURL, _ = url.Parse("/")
		}
	}

	if path.Base(redirectURL.Path) == "login" {
		redirectURL, _ = url.Parse("/")
		if redirectURL.Host != "" && redirectURL.Scheme != "" {
			redirectURL, _ = url.Parse(redirectURL.Scheme + "://" + redirectURL.Host)
		}
	}

	if auth.IsAuthenticated() {
		auth.context.Redirect(http.StatusSeeOther, redirectURL.String())
		return
	}

	session := sessions.Default(auth.context)
	session.Set("redirect", redirectURL.String())
	state := generateState()
	session.Set("state", state)
	sessionError := session.Save()
	if sessionError != nil {
		auth.context.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OIDC Login Session Error - %w", sessionError))
	}
	auth.context.Redirect(http.StatusTemporaryRedirect, auth.oauth2.AuthCodeURL(state))
}

func (auth *Authenticator) LogoutHandler() {
	if !auth.IsAuthenticated() {
		auth.context.Status(http.StatusOK)
		return
	}
	session := sessions.Default(auth.context)
	accessToken := session.Get("access_token")
	if accessToken != nil {
		logoutError := auth.revokeToken(accessToken.(string), "access_token")
		if logoutError != nil {
			auth.context.Error(fmt.Errorf("OIDC Revoke Access Token Error - %w", logoutError))
		}
	}
	refreshToken := session.Get("refresh_token")
	if refreshToken != nil {
		logoutError := auth.revokeToken(refreshToken.(string), "refresh_token")
		if logoutError != nil {
			auth.context.Error(fmt.Errorf("OIDC Revoke Refresh Token Error - %w", logoutError))
		}
	}
	session.Clear()
	session.Save()
	auth.context.Redirect(http.StatusTemporaryRedirect, auth.provider.Endpoints().EndSessionURL)
}

func (auth *Authenticator) RedirectHandler() {
	session := sessions.Default(auth.context)
	if auth.context.Query("state") != session.Get("state").(string) {
		auth.context.AbortWithError(http.StatusBadRequest, fmt.Errorf("OIDC error, state parameter is not valud"))
		return
	}

	token, tokenError := auth.oauth2.Exchange(oauth2.NoContext, auth.context.Query("code"))
	if tokenError != nil {
		auth.context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("OIDC error exchanging token, %w", tokenError))
		return
	}
	session.Set("access_token", token.AccessToken)
	session.Set("refresh_token", token.RefreshToken)

	IDTokenJWT, ok := token.Extra("id_token").(string)
	if !ok {
		auth.context.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OIDC error, ID token not found"))
		return
	} else {
		verifier := auth.provider.Verifier(&oidc.Config{ClientID: auth.oauth2.ClientID})
		_, tokenError := verifier.Verify(oauth2.NoContext, IDTokenJWT)
		if tokenError != nil {
			auth.context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("OIDC error verifying ID token, %w", tokenError))
			return
		}
		session.Set("id_token", IDTokenJWT)
	}

	userInfo, userError := auth.provider.UserInfo(oauth2.NoContext, auth.oauth2.TokenSource(oauth2.NoContext, token))
	if userError != nil {
		auth.context.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OIDC error fetching user information, %w", userError))
	} else {
		session.Set("user", userInfo.Subject)
	}

	if sessionError := session.Save(); sessionError != nil {
		auth.context.AbortWithError(http.StatusInternalServerError, fmt.Errorf("OIDC error storing session, %w", sessionError))
		return
	}

	redirectURL := session.Get("redirect").(string)
	auth.context.Redirect(http.StatusSeeOther, redirectURL)
}

func (auth *Authenticator) IsAuthenticated() bool {
	session := sessions.Default(auth.context)
	accessToken := session.Get("access_token")
	return accessToken != nil
}

func (auth *Authenticator) revokeToken(token string, tokenType string) error {

	client := http.DefaultClient

	request, HTTPError := http.NewRequest("POST", auth.provider.Endpoints().RevokeURL, nil)
	if HTTPError != nil {
		return HTTPError
	}

	request.SetBasicAuth(url.QueryEscape(auth.oauth2.ClientID), url.QueryEscape(auth.oauth2.ClientSecret))
	query := request.URL.Query()
	query.Add("token", token)
	query.Add("token_type_hint", tokenType)
	request.URL.RawQuery = query.Encode()
	response, HTTPError := client.Do(request)
	if HTTPError != nil {
		return HTTPError
	}
	defer response.Body.Close()
	return nil
}

func generateState() string {
	state := make([]byte, 32)
	_, stateError := rand.Read(state)
	if stateError != nil {
		panic(stateError)
	}
	return base64.StdEncoding.EncodeToString(state)
}
