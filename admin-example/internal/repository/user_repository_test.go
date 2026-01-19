package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql" // ドライバをインポート
	"github.com/stretchr/testify/assert"
)

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
	res, err := repo.Create(ctx, "テストユーザー")

	// 4. 検証
	assert.NoError(t, err)
	lastID, _ := res.LastInsertId()
	assert.True(t, lastID > 0)

	// 5. 後片付け（テストデータを消す、あるいはトランザクションでロールバックする）
	db.Exec("DELETE FROM users WHERE id = ?", lastID)
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
