package main

import (
	"database/sql"
	"fmt"
	"go-example/admin-example/internal/controller"
	"go-example/admin-example/internal/repository"
	"go-example/admin-example/internal/service"
	"io"

	_ "github.com/go-sql-driver/mysql"
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

	// 4. ルーティング
	e.GET("/users", ctrl.Index)
	e.POST("/users/:id/update", ctrl.Update)

	e.Logger.Fatal(e.Start(":8080"))
}
