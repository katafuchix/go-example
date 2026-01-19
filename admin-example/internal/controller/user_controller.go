package controller

import (
	"go-example/admin-example/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	svc service.UserService
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{svc: s}
}

// 一覧表示 (GET /users)
func (c *UserController) Index(ctx echo.Context) error {
	users, err := c.svc.GetList(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 2. テンプレートに渡すデータを map で定義
	// テンプレート側の {{range .Users}} とキー名を合わせるのがポイント
	data := map[string]interface{}{
		"Users": users,
	}

	// 3. レンダリング
	return ctx.Render(http.StatusOK, "users/index", data)
}

// 名前更新 (POST /users/:id/update)
func (c *UserController) Update(ctx echo.Context) error {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	name := ctx.FormValue("name")

	if err := c.svc.UpdateName(ctx.Request().Context(), id, name); err != nil {
		return err
	}
	return ctx.Redirect(http.StatusSeeOther, "/users")
}
