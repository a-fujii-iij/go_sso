package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OidcConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Issuer       string `json:"issuer"`
}

func main() {
	r := gin.Default()

	// 設定ファイルを読み込む
	configData, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	var config OidcConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// OIDCプロバイダーを設定
	provider, err := oidc.NewProvider(context.Background(), config.Issuer)
	if err != nil {
		log.Fatalf("failed to get provider: %v", err)
	}
	// OAuth2設定
	oauthConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// 認証用のエンドポイント
	r.GET("/auth", func(c *gin.Context) {
		state := "random" // 実際にはCSRF対策のためにランダムな値を生成する必要があります。
		url := oauthConfig.AuthCodeURL(state)
		c.Redirect(http.StatusFound, url)
	})

	// コールバック用のエンドポイント
	r.GET("/auth/callback", func(c *gin.Context) {
		// エラーハンドリングは省略しています。
		code := c.Query("code")
		oauth2Token, err := oauthConfig.Exchange(context.Background(), code)
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
		idToken, err := provider.Verifier(&oidc.Config{ClientID: config.ClientID}).Verify(context.Background(), rawIDToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID Token"})
			return
		}

		// IDトークンのクレームを取得
		var claims struct {
			Email string `json:"email"`
			// 他に必要なクレームがあればここに追加します。
		}
		if err := idToken.Claims(&claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get claims"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"email": claims.Email})
	})

	r.Run()

}
