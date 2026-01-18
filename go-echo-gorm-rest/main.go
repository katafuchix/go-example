package main

import (
	"go-example/go-echo-gorm-rest/controller"
	"go-example/go-echo-gorm-rest/model"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yosssi/ace"
)

// Ace用のレンダラー構造体
type AceRenderer struct{}

func (r *AceRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// 共通レイアウトと個別ページをガッチャンコしてコンパイル
	tpl, err := ace.Load("views/layout/base", "views/"+name, nil)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}

func connect(c echo.Context) error {
	db, _ := model.DB.DB()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		return c.String(http.StatusInternalServerError, "DB接続失敗しました")
	} else {
		return c.String(http.StatusOK, "DB接続しました")
	}
}

/*
	func main() {
		e := echo.New()
		e.GET("/", connect)
		e.Logger.Fatal(e.Start(":8080"))
	}
*/
func main() {
	e := echo.New()
	db, _ := model.DB.DB()
	defer db.Close()

	// 第一引数：ブラウザからアクセスする時のパス (/static/css/style.css など)
	// 第二引数：実際のサーバー上のディレクトリ名
	e.Static("/static", "public")

	// もしルート直下で画像などを配信したい場合
	// e.Static("/", "public")

	e.GET("/api/users", controller.GetUsers) // 追加
	e.GET("/users/:id", controller.GetUser)  // 追加
	e.POST("/users", controller.CreateUser)
	e.PUT("/users/:id", controller.UpdateUser)    // 追加
	e.DELETE("/users/:id", controller.DeleteUser) // 追加

	//e := echo.New()
	e.Renderer = &AceRenderer{} // レンダラーを登録

	e.GET("/users", func(c echo.Context) error {
		/*users := []map[string]interface{}{
			{"ID": 1, "Name": "田中"},
			{"ID": 2, "Name": "佐藤"},
		}*/
		users := []model.User{}

		// データベースから全てのユーザーを取得
		model.DB.Find(&users)

		// views/user/index.ace を呼び出す（baseと合体して表示される）
		return c.Render(200, "user/index", map[string]interface{}{
			"Users": users,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
