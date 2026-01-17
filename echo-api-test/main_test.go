// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

// 初期ユーザーデータ
var initialUsers = map[string]User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

// setupEcho はEchoのインスタンスを初期化し、ルートエンドポイントを登録します
func setupEcho() *echo.Echo {
	// 初期ユーザーデータで users マップをリセット
	users = make(map[string]User)
	for k, v := range initialUsers {
		users[k] = v
	}

	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())  // ログ記録ミドルウェア
	e.Use(middleware.Recover()) // パニック回復ミドルウェア
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// テスト用のベーシック認証設定
		if username == "admin" && password == "password" {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/hello", helloHandler)
	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.POST("/users", createUserHandler)
	e.PUT("/users/:id", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	return e
}

// 1. 基本的なGET /helloのテスト
func TestHelloHandler(t *testing.T) {
	e := setupEcho()

	// テスト用のGETリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	// レスポンスを記録するレコーダーを作成
	rec := httptest.NewRecorder()
	// Echoのコンテキストを作成
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, helloHandler(c)) {
		// ステータスコードが200であることを確認
		assert.Equal(t, http.StatusOK, rec.Code)
		// レスポンスボディが期待通りであることを確認
		assert.JSONEq(t, `{"message":"Hello, World!"}`, rec.Body.String())
	}
}

// 2. GET /users のテスト
func TestGetUsersHandler(t *testing.T) {
	e := setupEcho()

	// テスト用のGETリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, getUsersHandler(c)) {
		// ステータスコードが200であることを確認
		assert.Equal(t, http.StatusOK, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `[{"id":"1","name":"Alice"},{"id":"2","name":"Bob"}]`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 3. GET /users/:id の成功ケースのテスト
func TestGetUserHandler_Success(t *testing.T) {
	e := setupEcho()

	// 存在するユーザーIDでリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("1")

	// ハンドラを呼び出す
	if assert.NoError(t, getUserHandler(c)) {
		// ステータスコードが200であることを確認
		assert.Equal(t, http.StatusOK, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"id":"1","name":"Alice"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 4. GET /users/:id の失敗ケースのテスト（ユーザーが存在しない）
func TestGetUserHandler_NotFound(t *testing.T) {
	e := setupEcho()

	// 存在しないユーザーIDでリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("999")

	// ハンドラを呼び出す
	if assert.NoError(t, getUserHandler(c)) {
		// ステータスコードが404であることを確認
		assert.Equal(t, http.StatusNotFound, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"User not found"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 5. POST /users の成功ケースのテスト
func TestCreateUserHandler_Success(t *testing.T) {
	e := setupEcho()

	// 新しいユーザーを作成
	user := User{ID: "3", Name: "Charlie"}
	// ユーザーをJSONにエンコード
	body, _ := json.Marshal(user)
	// テスト用のPOSTリクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, createUserHandler(c)) {
		// ステータスコードが201であることを確認
		assert.Equal(t, http.StatusCreated, rec.Code)
		// レスポンスボディが期待通りであることを確認
		assert.JSONEq(t, `{"id":"3","name":"Charlie"}`, rec.Body.String())
		// ユーザーがマップに追加されていることを確認
		assert.Contains(t, users, "3")
	}
}

// 6. POST /users の失敗ケース（既に存在するユーザーID）のテスト
func TestCreateUserHandler_UserAlreadyExists(t *testing.T) {
	e := setupEcho()

	// 既に存在するユーザーIDでリクエストを作成
	user := User{ID: "1", Name: "Alice Duplicate"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, createUserHandler(c)) {
		// ステータスコードが400であることを確認
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"User already exists"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 7. POST /users の失敗ケース（無効なJSON）のテスト
func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	e := setupEcho()

	// 不正なJSONデータを作成
	invalidJSON := `{"id": "4", "name": }` // JSONの構文エラー
	// テスト用のPOSTリクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, createUserHandler(c)) {
		// ステータスコードが400であることを確認
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"Invalid input"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 8. POST /users の失敗ケース（フィールド不足）のテスト
func TestCreateUserHandler_MissingFields(t *testing.T) {
	e := setupEcho()

	// 必須フィールドが不足しているユーザーを作成
	user := User{ID: "", Name: "NoID"}
	// ユーザーをJSONにエンコード
	body, _ := json.Marshal(user)
	// テスト用のPOSTリクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラを呼び出す
	if assert.NoError(t, createUserHandler(c)) {
		// ステータスコードが400であることを確認
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"Missing fields"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 9. PUT /users/:id の成功ケースのテスト
func TestUpdateUserHandler_Success(t *testing.T) {
	e := setupEcho()

	// 既存のユーザーを更新
	updatedUser := User{Name: "Alice Updated"}
	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("1")

	// ハンドラを呼び出す
	if assert.NoError(t, updateUserHandler(c)) {
		// ステータスコードが200であることを確認
		assert.Equal(t, http.StatusOK, rec.Code)
		// レスポンスボディが期待通りであることを確認
		assert.JSONEq(t, `{"id":"1","name":"Alice Updated"}`, rec.Body.String())
		// ユーザーが更新されていることを確認
		assert.Equal(t, "Alice Updated", users["1"].Name)
	}
}

// 10. PUT /users/:id の失敗ケースのテスト（ユーザーが存在しない）
func TestUpdateUserHandler_NotFound(t *testing.T) {
	e := setupEcho()

	// 存在しないユーザーIDでリクエストを作成
	updatedUser := User{Name: "Nonexistent User"}
	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("999")

	// ハンドラを呼び出す
	if assert.NoError(t, updateUserHandler(c)) {
		// ステータスコードが404であることを確認
		assert.Equal(t, http.StatusNotFound, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"User not found"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 11. DELETE /users/:id の成功ケースのテスト
func TestDeleteUserHandler_Success(t *testing.T) {
	e := setupEcho()

	// 既存のユーザーを削除
	req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("2")

	// ハンドラを呼び出す
	if assert.NoError(t, deleteUserHandler(c)) {
		// ステータスコードが204であることを確認
		assert.Equal(t, http.StatusNoContent, rec.Code)
		// ユーザーがマップから削除されていることを確認
		assert.NotContains(t, users, "2")
	}
}

// 12. DELETE /users/:id の失敗ケースのテスト（ユーザーが存在しない）
func TestDeleteUserHandler_NotFound(t *testing.T) {
	e := setupEcho()

	// 存在しないユーザーIDでリクエストを作成
	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("999")

	// ハンドラを呼び出す
	if assert.NoError(t, deleteUserHandler(c)) {
		// ステータスコードが404であることを確認
		assert.Equal(t, http.StatusNotFound, rec.Code)
		// レスポンスボディが期待通りであることを確認
		expected := `{"error":"User not found"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

// 13. テーブル駆動テストの例（GET /hello と GET /users）
func TestHelloHandler_TableDriven(t *testing.T) {
	e := setupEcho()

	// テストケースの定義
	tests := []struct {
		name           string // テストケースの名前
		method         string // HTTPメソッド
		target         string // リクエストURL
		expectedStatus int    // 期待するステータスコード
		expectedBody   string // 期待するレスポンスボディ
	}{
		{
			name:           "Valid GET /hello",
			method:         http.MethodGet,
			target:         "/hello",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Hello, World!"}`,
		},
		{
			name:           "Valid GET /users",
			method:         http.MethodGet,
			target:         "/users",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","name":"Alice"},{"id":"2","name":"Bob"}]`,
		},
		// 他のエンドポイントのテストケースもここに追加可能
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のリクエストを作成
			var req *http.Request
			if tt.method == http.MethodGet {
				req = httptest.NewRequest(tt.method, tt.target, nil)
			} else {
				req = httptest.NewRequest(tt.method, tt.target, nil)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var handler echo.HandlerFunc
			// リクエストURLに応じてハンドラを選択
			switch tt.target {
			case "/hello":
				handler = helloHandler
			case "/users":
				handler = getUsersHandler
			// 他のエンドポイントのハンドラをここに追加
			default:
				t.Fatalf("Unknown target: %s", tt.target)
			}

			// ハンドラを呼び出す
			if assert.NoError(t, handler(c)) {
				// ステータスコードが期待通りであることを確認
				assert.Equal(t, tt.expectedStatus, rec.Code)
				// レスポンスボディが期待通りであることを確認
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

// 14. 認証が必要なエンドポイントへのアクセステスト
func TestAuthentication(t *testing.T) {
	e := setupEcho()

	// 正しい認証情報でのテスト
	t.Run("Valid Authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.SetBasicAuth("admin", "password")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// 誤った認証情報でのテスト
	t.Run("Invalid Authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.SetBasicAuth("admin", "wrongpassword")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// 認証情報なしでのテスト
	t.Run("No Authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
