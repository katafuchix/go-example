package main

import (
	"database/sql"
	"fmt"
	"go-example/admin-example/internal/controller"
	"go-example/admin-example/internal/repository"
	"go-example/admin-example/internal/service"
	"io"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yosssi/ace"
)

// Aceテンプレート用のRenderer設定
type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	path := "views/" + name
	fmt.Printf("Attempting to load template: %s.ace\n", path) // パスを表示

	//tpl, err := ace.Load("views/"+name, "", nil)
	// 第1引数: ベース(共通)テンプレートのパス
	// 第2引数: 中身(Indexなど)のテンプレートのパス
	tpl, err := ace.Load("views/layout/base", "views/"+name, nil)
	if err != nil {
		// どこでエラーが起きているかターミナルに出力する
		fmt.Printf("Ace Load Error: %v\n", err)
		return err
	}
	return tpl.Execute(w, data)
}

// CustomValidator はEchoのValidatorインターフェースを満たす構造体
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// 1. DB接続 (ここは共通パッケージ等に切り出してもOK)
	db, _ := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test?parseTime=true")

	// 2. DI (依存性の注入)
	repo := repository.NewUserRepository(db)  // Repoを作る
	svc := service.NewUserService(repo)       // RepoをServiceに入れる
	ctrl := controller.NewUserController(svc) // ServiceをControllerに入れる

	// 3. Echoの起動
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = &TemplateRenderer{}

	// 1. セッションの設定を追加（これが今回のエラーの直接の原因）
	// "secret-key" は適当な文字列でOKです。これがセッションの暗号化に使われます。
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))

	// バリデーターを登録
	e.Validator = &CustomValidator{validator: validator.New()}

	// CSRFミドルウェアを登録
	/*e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf", // フォーム内の <input name="csrf"> を見る
	}))*/

	// 静的ファイルの配信
	e.Static("/static", "public")

	// 4. ルーティング
	e.GET("/users", ctrl.Index)
	e.POST("/users/:id/update", ctrl.Update)

	// 【新規登録】
	// 1. 画面を表示する
	e.GET("/users/create", ctrl.New) // ← これを足す
	// 2. フォームから送られてきたデータを保存する
	e.POST("/users", ctrl.Create)

	// 【編集・更新】
	// 1. 画面を表示する (IDをURLに含む)
	e.GET("/users/edit/:id", ctrl.Edit) // ← これを足す
	// 2. データを更新する
	e.POST("/users/update", ctrl.Update)

	e.Logger.Fatal(e.Start(":8080"))
}
