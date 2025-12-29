package user

import (
	"context"
	"time"
)

// MockRepository はRepositoryのモック実装
// テストと開発用に定数値を返す
type MockRepository struct{}

// NewMockRepository は新しいモックユーザーリポジトリを作成する
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

// FindByAuth0UserID は定数値を持つモックユーザーを返す
func (r *MockRepository) FindByAuth0UserID(ctx context.Context, auth0UserID string) (*User, error) {
	// モックデータを返す
	now := time.Now()
	return &User{
		ID:              "llu_mock_user_001",
		Auth0UserID:     auth0UserID,
		WorkspaceID:     "ws_mock_workspace_001",
		IsPrivileged:    false,
		IdPConnectionID: nil,
		Email:           "user01@example.com",
		Name:            "Mock User 01",
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Update はモック実装で何もしない
func (r *MockRepository) Update(ctx context.Context, user *User) error {
	// モック実装 - 実際のデータベース更新は行わない
	// タイムスタンプのみ更新
	user.UpdatedAt = time.Now()
	return nil
}
