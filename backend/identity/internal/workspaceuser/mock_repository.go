package workspaceuser

import (
	"context"
	"fmt"
	"time"
)

// MockRepository はWorkspaceUserのモックリポジトリ
type MockRepository struct {
	users map[string]*WorkspaceUser
}

// NewMockRepository は新しいモックリポジトリを作成する
func NewMockRepository() *MockRepository {
	// モックデータを初期化
	users := map[string]*WorkspaceUser{
		"auth0|6952b421821fed371daac9df": {
			ID:          "wsu-001",
			WorkspaceID: "ws-001",
			Auth0UserID: "auth0|6952b421821fed371daac9df",
			Email:       "user01@example.com",
			Name:        "User 01",
			CreatedAt:   time.Now().Add(-10 * 24 * time.Hour), // 10日前
		},
	}

	return &MockRepository{
		users: users,
	}
}

// FindByAuth0UserID はAuth0ユーザーIDでWorkspaceUserを取得する
func (r *MockRepository) FindByAuth0UserID(ctx context.Context, auth0UserID string) (*WorkspaceUser, error) {
	user, ok := r.users[auth0UserID]
	if !ok {
		return nil, fmt.Errorf("workspace user not found: %s", auth0UserID)
	}
	return user, nil
}
