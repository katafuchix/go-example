package controller

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// すべてのコントローラーの親になる構造体
type BaseController struct{}

// 二重送信防止：トークンチェック
func (b *BaseController) IsValidAndDestroyToken(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	saved := sess.Values["submit_token"]
	submitted := c.FormValue("csrf")

	if saved == nil || submitted != saved.(string) {
		return false
	}
	delete(sess.Values, "submit_token")
	sess.Save(c.Request(), c.Response())
	return true
}

// トークン発行
func (b *BaseController) IssueToken(c echo.Context) string {
	token := uuid.NewString()
	sess, _ := session.Get("session", c)
	sess.Values["submit_token"] = token
	sess.Save(c.Request(), c.Response())
	return token
}

func (b *BaseController) GetValidationErrMsg(err error, form interface{}) string {
	if err == nil {
		return ""
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return "入力エラーです"
	}

	errField := ve[0]

	// 型情報を取得
	t := reflect.TypeOf(form)
	// もしポインタ (*UserCreateForm) が渡されたら、その中身 (UserCreateForm) に移動する
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// フィールド名で検索してタグを取得
	field, found := t.FieldByName(errField.Field())
	if !found {
		return errField.Field() + "の入力が正しくありません"
	}

	label := field.Tag.Get("label")
	if label == "" {
		label = errField.Field() // labelタグがない場合はフィールド名
	}

	// ルールごとに日本語化
	switch errField.Tag() {
	case "required":
		return label + "を入力してください"
	case "min":
		return label + "は" + errField.Param() + "文字以上で入力してください"
	case "max":
		return label + "は" + errField.Param() + "文字以内で入力してください"
	default:
		return label + "が正しくありません"
	}
}
