package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// infraパッケージのRandString関数をモックするためのインターフェース
type MockInfra struct {
	mock.Mock
}

func (m *MockInfra) RandString(n int) string {
	args := m.Called(n)
	return args.String(0)
}

// oauth2.ConfigのAuthCodeURL関数をモックするためのインターフェース
type MockOAuth2Config struct {
	mock.Mock
}

func (m *MockOAuth2Config) AuthCodeURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func TestHandleAuthRequest(t *testing.T) {
	// モックオブジェクトの作成
	mockInfra := new(MockInfra)
	mockOAuth2Config := new(MockOAuth2Config)

	// モックの期待値を設定
	mockInfra.On("RandString", 16).Return("mocked_state")
	mockOAuth2Config.On("AuthCodeURL", "mocked_state").Return("mocked_url")

	// OIDCHandlerのインスタンスを作成
	handler := OIDCHandler{
		oauthConfig: mockOAuth2Config,
	}

	// gin.Contextのモックを作成
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// テスト対象の関数を実行
	handler.HandleAuthRequest(c)

	// リダイレクトが適切に行われたか検証
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "mocked_url", w.Header().Get("Location"))

	// モックの呼び出しを検証
	mockInfra.AssertExpectations(t)
	mockOAuth2Config.AssertExpectations(t)
}
