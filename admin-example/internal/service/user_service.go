package service

import (
	"context"
	"go-example/admin-example/internal/repository"
)

// サービス側のインターフェース（Controllerがこれを使う）
type UserService interface {
	Register(ctx context.Context, name string) error
	GetList(ctx context.Context) ([]repository.User, error)
	UpdateName(ctx context.Context, id uint64, name string) error
	DeleteUser(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (repository.User, error) // これを追加
}

type userService struct {
	repo repository.UserRepository // Repositoryの「インターフェース」を持つ
}

// コンストラクタ（ここでRepositoryの実体を注入する）
func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

// --- 実装 ---

func (s *userService) Register(ctx context.Context, name string) error {
	// 必要ならここに「名前の重複チェック」などのバリデーションロジックを書く
	_, err := s.repo.Create(ctx, name)
	return err
}

func (s *userService) GetList(ctx context.Context) ([]repository.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) UpdateName(ctx context.Context, id uint64, name string) error {
	return s.repo.Update(ctx, id, name)
}

func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) FindByID(ctx context.Context, id uint64) (repository.User, error) {
	// Repositoryの FindByID を呼び出すだけ
	return s.repo.FindByID(ctx, id)
}
