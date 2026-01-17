// main.go
package main

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// User はユーザー情報を表す構造体です
type User struct {
	ID   string `json:"id"`   // ユーザーID
	Name string `json:"name"` // ユーザー名
}

// サンプルユーザーデータを保持するマップ
var users = map[string]User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

func main() {
	// Echoのインスタンスを作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())  // ログ記録ミドルウェア
	e.Use(middleware.Recover()) // パニック回復ミドルウェア
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// シンプルなベーシック認証の例
		// 実際のプロジェクトでは、より安全な認証方法を使用してください
		if username == "admin" && password == "password" {
			return true, nil
		}
		return false, nil
	}))

	// ルートエンドポイントとハンドラを登録
	e.GET("/hello", helloHandler)             // GET /hello
	e.GET("/users", getUsersHandler)          // GET /users
	e.GET("/users/:id", getUserHandler)       // GET /users/:id
	e.POST("/users", createUserHandler)       // POST /users
	e.PUT("/users/:id", updateUserHandler)    // PUT /users/:id
	e.DELETE("/users/:id", deleteUserHandler) // DELETE /users/:id

	// サーバーをポート8080で開始
	e.Start(":8080")
}

// helloHandler は /hello エンドポイントのハンドラです
func helloHandler(c echo.Context) error {
	// JSON形式でメッセージを返す
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

// getUsersHandler は /users エンドポイントのハンドラです
func getUsersHandler(c echo.Context) error {
	// ユーザーリストをスライスに変換
	userList := []User{}
	for _, user := range users {
		userList = append(userList, user)
	}

	// ユーザーリストをIDでソート
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].ID < userList[j].ID
	})

	// JSON形式でソートされたユーザーリストを返す
	return c.JSON(http.StatusOK, userList)
}

// getUserHandler は /users/:id エンドポイントのハンドラです
func getUserHandler(c echo.Context) error {
	// URLパラメータからユーザーIDを取得
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		// ユーザーが存在しない場合は404を返す
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}
	// ユーザー情報をJSON形式で返す
	return c.JSON(http.StatusOK, user)
}

// createUserHandler は /users エンドポイントのハンドラです
func createUserHandler(c echo.Context) error {
	user := new(User)
	// リクエストボディをUser構造体にバインド
	if err := c.Bind(user); err != nil {
		// バインドに失敗した場合は400を返す
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}
	// 必須フィールドのチェック
	if user.ID == "" || user.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing fields",
		})
	}
	// ユーザーが既に存在するか確認
	if _, exists := users[user.ID]; exists {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User already exists",
		})
	}
	// ユーザーをマップに追加
	users[user.ID] = *user
	// 作成したユーザー情報を201で返す
	return c.JSON(http.StatusCreated, user)
}

// updateUserHandler は /users/:id エンドポイントのハンドラです
func updateUserHandler(c echo.Context) error {
	// URLパラメータからユーザーIDを取得
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		// ユーザーが存在しない場合は404を返す
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	updatedUser := new(User)
	// リクエストボディをUser構造体にバインド
	if err := c.Bind(updatedUser); err != nil {
		// バインドに失敗した場合は400を返す
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}
	// 必須フィールドのチェック
	if updatedUser.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing fields",
		})
	}

	// ユーザー情報を更新
	user.Name = updatedUser.Name
	users[id] = user

	// 更新したユーザー情報を200で返す
	return c.JSON(http.StatusOK, user)
}

// deleteUserHandler は /users/:id エンドポイントのハンドラです
func deleteUserHandler(c echo.Context) error {
	// URLパラメータからユーザーIDを取得
	id := c.Param("id")
	_, exists := users[id]
	if !exists {
		// ユーザーが存在しない場合は404を返す
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}
	// ユーザーをマップから削除
	delete(users, id)
	// 削除成功を204で返す（ボディなし）
	return c.NoContent(http.StatusNoContent)
}
