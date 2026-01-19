package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql" // ドライバをインポート
	"github.com/stretchr/testify/assert"
)

/*
func TestUserRepository_Create(t *testing.T) {
	// 1. テスト用DBに接続（本番用とは別のDBを用意すること！）
	db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/test_db?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// 2. リポジトリの準備（sqlcが生成したNew関数を使用）
	repo := New(db)
	ctx := context.Background()

	// 3. データの挿入（テスト実行）
	res, err := repo.CreateUser(ctx, sql.NullString{
		String: "テストユーザー",
		Valid:  true, // これをtrueにしないと、DBにはNULLとして保存されてしまいます
	})

	// 4. 検証
	assert.NoError(t, err)
	lastID, _ := res.LastInsertId()
	assert.True(t, lastID > 0)

	// 5. 後片付け（テストデータを消す、あるいはトランザクションでロールバックする）
	db.Exec("DELETE FROM users WHERE id = ?", lastID)
}
*/

func TestUserRepository_Create(t *testing.T) {
	// 1. 本物のDB接続
	db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/test_db?parseTime=true")
	defer db.Close()

	// 2. ★ここが最重要：あなたが作った「NewUserRepository」を呼ぶ
	// これで repo は UserRepository インターフェース（実体は *userRepository）になる
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 3. インターフェースで定義した「Create」を呼ぶ
	// 引数は string で渡せる（内部で sql.NullString に変換されるのをテストできる）
	res, err := repo.Create(ctx, "テストユーザー")

	// 4. 検証
	assert.NoError(t, err)
	lastID, _ := res.LastInsertId()
	assert.True(t, lastID > 0)

	// 後片付け
	db.Exec("DELETE FROM users WHERE id = ?", lastID)
}

func TestUserRepository(t *testing.T) {
	// 1. テストDB準備
	db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/test_db?parseTime=true")
	defer db.Close()

	// 2. 「あなたの作った」ラッパーを生成
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Createのテスト", func(t *testing.T) {
		// ここではインターフェースのメソッド名「Create」を呼ぶ
		// 引数も sql.NullString ではなく、ただの string でOK（ラッパーが変換してくれるから）
		res, err := repo.Create(ctx, "テスト太郎")

		assert.NoError(t, err)
		lastID, _ := res.LastInsertId()
		assert.True(t, lastID > 0)

		// 後片付け
		db.Exec("DELETE FROM users WHERE id = ?", lastID)
	})

	t.Run("Listのテスト", func(t *testing.T) {
		users, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
}

/*
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB() // テスト用DB接続

    // トランザクション開始
    tx, _ := db.Begin()
    defer tx.Rollback() // テストが終わったら何があっても元に戻す！

    // sqlcの生成したクエリに tx (トランザクション) を渡す
    repo := New(tx)

    // 実行 & 検証
    _, err := repo.Create(context.Background(), "テストユーザー")
    assert.NoError(t, err)

    // DBには何も残らないので、次のテストも真っさらな状態で始められる
}
*/
