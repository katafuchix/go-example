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

// 複数のエラーをまとめて取得する
func (b *BaseController) GetValidationErrors(err error, form interface{}) map[string]string {
	errors := make(map[string]string)
	if err == nil {
		return errors
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	t := reflect.TypeOf(form)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for _, fe := range ve {
		// フィールド名とラベルの取得
		label := fe.Field()
		if field, found := t.FieldByName(fe.Field()); found {
			l := field.Tag.Get("label")
			if l != "" {
				label = l
			}
		}

		// エラー内容に応じたメッセージ作成
		var msg string
		switch fe.Tag() {
		case "required":
			msg = label + "を入力してください"
		case "min":
			msg = label + "は" + fe.Param() + "文字以上で入力してください"
		case "max":
			msg = label + "は" + fe.Param() + "文字以内で入力してください"
		default:
			msg = label + "が正しくありません"
		}

		// フィールド名（小文字のnameなど）をキーにして保存
		errors[fe.Field()] = msg
	}
	return errors
}
