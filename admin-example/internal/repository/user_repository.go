package repository

import (
	"context"
	"database/sql"
)

// 1. インターフェース定義（テストでMockに差し替えるための「契約」）
type UserRepository interface {
	Create(ctx context.Context, name string) (sql.Result, error)
	FindByID(ctx context.Context, id uint64) (User, error)
	Update(ctx context.Context, id uint64, name string) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context) ([]User, error)
}

// 2. 実体となる構造体
type userRepository struct {
	q  *Queries
	db *sql.DB
}

// 3. コンストラクタ（sqlcの生成物をラップして返す）
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		q:  New(db), // Newはsqlcが生成した関数
		db: db,
	}
}

// --- 以下、インターフェースを満たすための実装 ---

func (r *userRepository) Create(ctx context.Context, name string) (sql.Result, error) {
	return r.q.CreateUser(ctx, sql.NullString{String: name, Valid: name != ""})
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (User, error) {
	return r.q.GetUser(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, id uint64, name string) error {
	return r.q.UpdateUser(ctx, UpdateUserParams{
		ID:   id,
		Name: sql.NullString{String: name, Valid: name != ""},
	})
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context) ([]User, error) {
	return r.q.ListUsers(ctx)
}
