package main

import (
	"go-example/go-echo-gorm-rest/controller"
	"go-example/go-echo-gorm-rest/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

	e.GET("/users", controller.GetUsers)    // 追加
	e.GET("/users/:id", controller.GetUser) // 追加
	e.POST("/users", controller.CreateUser)
	e.PUT("/users/:id", controller.UpdateUser)    // 追加
	e.DELETE("/users/:id", controller.DeleteUser) // 追加
	e.Logger.Fatal(e.Start(":8080"))
}
