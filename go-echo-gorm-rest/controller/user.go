package controller

import (
	"net/http"

	"go-example/go-echo-gorm-rest/model"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	// 新しいユーザーを表す構造体を宣言
	user := model.User{}

	// コンテキストからユーザーデータをバインド
	if err := c.Bind(&user); err != nil {
		// バインドに失敗した場合、エラーを返す
		return err
	}

	// データベースに新しいユーザーを作成
	model.DB.Create(&user)

	// 作成されたユーザーをJSON形式で返す
	return c.JSON(http.StatusCreated, user)
}

func GetUsers(c echo.Context) error {
	// ユーザーのスライスを宣言
	users := []model.User{}

	// データベースから全てのユーザーを取得
	model.DB.Find(&users)

	// 取得したユーザーをJSON形式で返す
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	// 単一のユーザーを表す構造体を宣言
	user := model.User{}

	// コンテキストからユーザーデータをバインド
	if err := c.Bind(&user); err != nil {
		// バインドに失敗した場合、エラーを返す
		return err
	}

	// データベースから単一のユーザーを取得
	model.DB.Take(&user)

	// 取得したユーザーをJSON形式で返す
	return c.JSON(http.StatusOK, user)
}

// $ curl -X POST -H "Content-Type: application/json" -d '{"name":"テスト太郎"}' localhost:8080/users
//{"id":1,"name":"テスト太郎","created_at":"2022-12-30T08:41:51.838Z","updated_at":"2022-12-30T08:41:51.838Z"}

func UpdateUser(c echo.Context) error {
	// 0. URLからIDを取得する（ブラウザからのID指定を無視する）
	id := c.Param("id") // URLからIDを取得
	user := model.User{}

	// 1. 対象のユーザーをIDで検索
	if err := model.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// 2. リクエストボディの内容を構造体にバインド（上書き）
	// ここで Name などを書き換える
	// c.Bind は送られてきた項目だけを上書きし、送られなかった項目（パスワード等）は維持する
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// 3. データベースに保存
	// ここで GORM が自動的に UpdatedAt を現在時刻に更新します
	model.DB.Save(&user)

	return c.JSON(http.StatusOK, user)
}

//  curl -X PUT -H "Content-Type: application/json" -d '{"name":"山田たろう"}' localhost:8080/users/1
// {"id":1,"name":"山田たろう","created_at":"2026-01-18T10:56:54.381+09:00","updated_at":"2026-01-18T11:16:06.099+09:00"}

/*
func UpdateUserName(c echo.Context) error {
    id := c.Param("id")
    var user model.User

    if err := c.Bind(&user); err != nil {
        return err
    }

    // IDをURLから取得したものに強制的に書き換える（改ざん防止）
    uid, _ := strconv.Atoi(id)
    user.ID = uint(uid)

    // Nameカラムだけを更新対象にする。
    // もしリクエストに "is_admin: true" とか入っていても無視されるので安全！
    model.DB.Model(&user).Select("Name").Updates(user)

    return c.JSON(http.StatusOK, user)
}
*/

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	user := model.User{}

	// IDを指定して削除を実行
	// model.User{} を渡すことで、どのテーブルから消すかをGORMが判断します
	if err := model.DB.Delete(&user, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Delete failed")
	}

	// 204 No Content を返すのが一般的です
	return c.NoContent(http.StatusNoContent)
}
