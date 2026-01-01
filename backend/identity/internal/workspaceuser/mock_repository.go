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
		"auth0|user002": {
			ID:          "wsu-002",
			WorkspaceID: "ws-001",
			Auth0UserID: "auth0|user002",
			Email:       "user02@example.com",
			Name:        "User 02",
			CreatedAt:   time.Now().Add(-9 * 24 * time.Hour), // 9日前
		},
		"auth0|user003": {
			ID:          "wsu-003",
			WorkspaceID: "ws-001",
			Auth0UserID: "auth0|user003",
			Email:       "user03@example.com",
			Name:        "User 03",
			CreatedAt:   time.Now().Add(-8 * 24 * time.Hour), // 8日前
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

// ListByWorkspaceID はワークスペースIDでWorkspaceUserの一覧を取得する
func (r *MockRepository) ListByWorkspaceID(ctx context.Context, workspaceID string, pageSize int32, pageToken string) ([]*WorkspaceUser, string, error) {
	// ワークスペースIDに一致するユーザーを収集
	var users []*WorkspaceUser
	for _, user := range r.users {
		if user.WorkspaceID == workspaceID {
			users = append(users, user)
		}
	}

	// CreatedAt の降順でソート
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			if users[i].CreatedAt.Before(users[j].CreatedAt) {
				users[i], users[j] = users[j], users[i]
			}
		}
	}

	// ページングを簡易実装（実際のDBではoffsetやcursorを使用）
	if pageSize == 0 {
		pageSize = 10 // デフォルト
	}

	var result []*WorkspaceUser
	var nextPageToken string

	// 簡易的にページトークンを数値として扱う
	start := 0
	if pageToken != "" {
		fmt.Sscanf(pageToken, "%d", &start)
	}

	end := start + int(pageSize)
	if end > len(users) {
		end = len(users)
	}

	if start < len(users) {
		result = users[start:end]
		if end < len(users) {
			nextPageToken = fmt.Sprintf("%d", end)
		}
	}

	return result, nextPageToken, nil
}
