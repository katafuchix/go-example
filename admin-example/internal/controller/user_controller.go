package controller

import (
	"go-example/admin-example/internal/model"
	"go-example/admin-example/internal/service"
	"net/http"
	"strconv"

	// Echo v4に対応したものを呼び出す
	"github.com/labstack/echo/v4"
)

type UserController struct {
	BaseController // ★これを書くだけで Base の機能がすべて使える
	svc            service.UserService
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

/*
// 名前更新 (POST /users/:id/update)
func (c *UserController) Update(ctx echo.Context) error {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	name := ctx.FormValue("name")

	if err := c.svc.UpdateName(ctx.Request().Context(), id, name); err != nil {
		return err
	}
	return ctx.Redirect(http.StatusSeeOther, "/users")
}
*/

// New: 登録画面を表示するだけ
func (u *UserController) New(c echo.Context) error {
	// 新規画面表示時にトークンを発行
	return c.Render(http.StatusOK, "users/new", map[string]interface{}{
		"csrf": u.IssueToken(c),
	})
}
func (u *UserController) Create(c echo.Context) error {
	// 1. 二重送信チェック (バックボタン対策)
	if !u.IsValidAndDestroyToken(c) {
		return c.String(http.StatusBadRequest, "二重送信エラーです。前の画面に戻ってやり直してください。")
	}

	form := new(model.UserCreateForm)
	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		// ここ！ form がポインタなので、BaseController側で「label」タグを読み取れます
		msg := u.GetValidationErrMsg(err, form)
		return u.renderNew(c, msg, form.Name)
	}

	if err := u.svc.Register(c.Request().Context(), form.Name); err != nil {
		return u.renderNew(c, "保存に失敗しました", form.Name)
	}

	return c.Redirect(http.StatusSeeOther, "/users")
}

func (u *UserController) Edit(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := u.svc.FindByID(c.Request().Context(), id)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/users")
	}

	return c.Render(http.StatusOK, "users/edit", map[string]interface{}{
		"ID":   user.ID,
		"Name": user.Name.String,
		"csrf": u.IssueToken(c), // 編集画面表示時もトークン発行
	})
}

func (u *UserController) Update(c echo.Context) error {
	// 1. 二重送信チェック
	if !u.IsValidAndDestroyToken(c) {
		return c.String(http.StatusBadRequest, "二重送信エラーです。")
	}

	form := new(model.UserUpdateForm)
	if err := c.Bind(form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		// 共通バリデーションメッセージ関数を呼び出し
		msg := u.GetValidationErrMsg(err, form)
		return u.renderNew(c, msg, form.Name)
	}

	if err := u.svc.UpdateName(c.Request().Context(), form.ID, form.Name); err != nil {
		return u.renderEdit(c, form.ID, "更新に失敗しました", form.Name)
	}

	return c.Redirect(http.StatusSeeOther, "/users")
}

// --- ヘルパーメソッド (エラー時に新しいトークンを付けて再表示) ---

func (u *UserController) renderNew(c echo.Context, errMsg string, name string) error {
	return c.Render(http.StatusOK, "users/new", map[string]interface{}{
		"Error": errMsg,
		"Name":  name,
		"csrf":  u.IssueToken(c), // 再表示のたびにトークンを更新
	})
}

func (u *UserController) renderEdit(c echo.Context, id uint64, errMsg string, name string) error {
	return c.Render(http.StatusOK, "users/edit", map[string]interface{}{
		"ID":    id,
		"Error": errMsg,
		"Name":  name,
		"csrf":  u.IssueToken(c), // 再表示のたびにトークンを更新
	})
}
