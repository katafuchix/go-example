package service

import (
	"context"
	"database/sql"
	"testing"

	"go-example/admin-example/internal/repository"

	"github.com/stretchr/testify/assert"
)

// % go test -v
// % go test -cover
// % go test -coverprofile=cover.out
// % go tool cover -html=cover.out

// 1. Repositoryの偽物（Mock）を定義
// repository.UserRepository インターフェースを満たすように作る
type mockUserRepository struct {
	// 埋め込みなどはせず、必要なメソッドだけ定義してもOKですが、
	// 全てのメソッドを定義する必要があります。
}

/*
// 埋め込みを使った時短テクニック
type mockUserRepository struct {
    repository.UserRepository // これを書くと、実装していないメソッドがあってもコンパイルが通る
}
*/

func (m *mockUserRepository) List(ctx context.Context) ([]repository.User, error) {
	// テスト用のダミーデータを返す
	return []repository.User{
		{ID: 1, Name: sql.NullString{String: "テスト太郎", Valid: true}},
		{ID: 2, Name: sql.NullString{String: "テスト次郎", Valid: true}},
	}, nil
}

// --- 以下、インターフェースを満たすためのダミー実装 ---
// sqlcの生成するインターフェースに合わせて戻り値を (sql.Result, error) に統一します

func (m *mockUserRepository) Create(ctx context.Context, name string) (sql.Result, error) {
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, id uint64, name string) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

// 必要に応じて、sqlcが生成した他のメソッド（GetUserなど）があればここに追加します
func (m *mockUserRepository) GetUser(ctx context.Context, id uint64) (repository.User, error) {
	return repository.User{}, nil
}

// FindByID を追加（中身は空でOK）
func (m *mockUserRepository) FindByID(ctx context.Context, id uint64) (repository.User, error) {
	return repository.User{}, nil
}

// 2. 実際のテスト関数
func TestUserService_GetList(t *testing.T) {
	// A. 準備: 偽物のRepoを作り、Serviceに注入(DI)する
	mockRepo := &mockUserRepository{}
	svc := NewUserService(mockRepo)

	// B. 実行: Serviceのメソッドを呼ぶ
	users, err := svc.GetList(context.Background())

	// C. 検証: 結果が期待通りかチェックする
	assert.NoError(t, err)                         // エラーが出ていないこと
	assert.Len(t, users, 2)                        // 2件取得できていること
	assert.Equal(t, "テスト太郎", users[0].Name.String) // 1件目の名前が正しいこと
}
