package workspace

import (
	"context"
	"fmt"
	"time"
)

// MockRepository はWorkspaceのモックリポジトリ
type MockRepository struct {
	workspaces map[string]*Workspace
}

// NewMockRepository は新しいモックリポジトリを作成する
func NewMockRepository() *MockRepository {
	// モックデータを初期化
	workspaces := map[string]*Workspace{
		"ws-001": {
			ID:        "ws-001",
			Name:      "My Workspace",
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour), // 30日前
		},
	}

	return &MockRepository{
		workspaces: workspaces,
	}
}

// FindByID はIDでWorkspaceを取得する
func (r *MockRepository) FindByID(ctx context.Context, id string) (*Workspace, error) {
	workspace, ok := r.workspaces[id]
	if !ok {
		return nil, fmt.Errorf("workspace not found: %s", id)
	}
	return workspace, nil
}
