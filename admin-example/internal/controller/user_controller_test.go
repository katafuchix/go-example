package controller

import (
	"context"
	"database/sql"
	"go-example/admin-example/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// 1. Serviceの偽物（Mock）を定義
type mockUserService struct{}

func (m *mockUserService) GetList(ctx context.Context) ([]repository.User, error) {
	// サービスが返すダミーデータ
	return []repository.User{
		{ID: 1, Name: sql.NullString{String: "コントローラーテスト", Valid: true}},
	}, nil
}

// 他のメソッドもインターフェースを満たすために定義
func (m *mockUserService) Register(ctx context.Context, name string) error              { return nil }
func (m *mockUserService) UpdateName(ctx context.Context, id uint64, name string) error { return nil }
func (m *mockUserService) DeleteUser(ctx context.Context, id uint64) error              { return nil }

// 2. Rendererの偽物（Mock） ★ここがポイント
type mockRenderer struct{}

func (m *mockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// テンプレートファイルを読み込まず、渡されたデータの中身を強引にテキストとして書き出す
	// これにより「描画エラー」を防ぎつつ、データの受け渡しが正しいか検証できる
	if d, ok := data.(map[string]interface{}); ok {
		if users, ok := d["Users"].([]repository.User); ok && len(users) > 0 {
			w.Write([]byte(users[0].Name.String))
		}
	}
	return nil
}

func TestUserController_Index(t *testing.T) {
	// 2. Echoのセットアップ
	e := echo.New()

	// 重要：テスト中にAceテンプレートをレンダリングできるように設定
	// 本番と同じパスを指定するか、RendererをMockにする必要があります。
	// ここでは簡易的にRendererが設定されている前提とします。
	// e.Renderer = ... (main.goでの設定と同様のもの)

	// Rendererの偽物を登録
	e.Renderer = &mockRenderer{}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 3. DI：MockServiceを注入してControllerを作成
	mockSvc := &mockUserService{}
	ctrl := NewUserController(mockSvc)

	// 4. 実行
	// もしRendererの設定が難しい場合は、Index内で呼び出している
	// c.Render をモックアウトするか、テスト用のテンプレートパスを通します。
	err := ctrl.Index(c)

	// 5. 検証
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// HTMLの中身に、Mockが返した名前が含まれているかチェック
	assert.Contains(t, rec.Body.String(), "コントローラーテスト")
}

// 1. 常にエラーを返すServiceの偽物
type errorUserService struct {
	mockUserService // 他のメソッドを共通化
}

func (m *errorUserService) GetList(ctx context.Context) ([]repository.User, error) {
	return nil, sql.ErrConnDone // わざと「DB接続切れ」エラーを返す
}

func TestUserController_Index_Error(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 2. エラーを出すServiceを注入
	mockSvc := &errorUserService{}
	ctrl := NewUserController(mockSvc)

	// 3. 実行
	err := ctrl.Index(c)

	// 4. 検証
	// 通常、Echoのハンドラーがエラーを返すと、err に値が入ります
	assert.Error(t, err)

	// もしController内で HTTPError に変換しているなら、その中身をチェック
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusInternalServerError, he.Code)
	}
}
