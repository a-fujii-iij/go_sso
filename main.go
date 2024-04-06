package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"

	"github.com/a-fujii-iij/go_sso_gin/handler"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OidcConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Issuer       string `json:"issuer"`
}

const rs3Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs3Letters[int(rand.Int63()%int64(len(rs3Letters)))]
	}
	return string(b)
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

	oidcHandler := handler.NewOIDCHandler(oauthConfig, provider)

	// 認証用のエンドポイント
	r.GET("/auth", oidcHandler.HandleAuthRequest)

	// コールバック用のエンドポイント
	r.GET("/auth/callback", oidcHandler.HandleAuthResponse)

	r.Run()

}
