package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/a-fujii-iij/go_sso_gin/infra"
)

type OIDCHandler struct {
	oauthConfig oauth2.Config
	provider    *oidc.Provider
}

func NewOIDCHandler(oauthConfig oauth2.Config, provider *oidc.Provider) *OIDCHandler {
	return &OIDCHandler{oauthConfig: oauthConfig, provider: provider}
}

func (h *OIDCHandler) HandleAuthRequest(c *gin.Context) {
	state := infra.RandString(16)
	url := h.oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *OIDCHandler) HandleAuthResponse(c *gin.Context) {
	code := c.Query("code")
	oauth2Token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// IDトークンの取得
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	// IDトークンの検証
	idToken, err := h.provider.Verifier(&oidc.Config{ClientID: h.oauthConfig.ClientID}).Verify(context.Background(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID Token"})
		return
	}

	// IDトークンのクレームを取得
	var claims struct {
		Email string   `json:"email"`
		Sub   string   `json:"sub"`
		Aud   []string `json:"aud"`
		Iss   string   `json:"iss"`
		// 他に必要なクレームがあればここに追加します。
	}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("idToken.Claims failed. err: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clames": claims})
}
