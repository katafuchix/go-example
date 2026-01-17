package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Template はEchoに登録するためのレンダラー構造体です
type Template struct {
	templates *template.Template
}

// Render はEchoから呼び出される描画メソッドです
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echoのインスタンスを作成
	e := echo.New()

	// テンプレートの読み込み（viewsフォルダ内の.htmlファイルを全部読み込む設定）
	renderer := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = renderer

	// ルーティング： "/" にアクセスしたら文字列を返す
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Echoサーバー起動成功！")
	})

	e.GET("/hello", func(c echo.Context) error {
		// "hello.html" という名前のテンプレートを、データ（名前）付きで表示
		return c.Render(http.StatusOK, "hello.html", "管理者")
	})

	// 8080ポートで起動
	e.Logger.Fatal(e.Start(":8080"))
}
