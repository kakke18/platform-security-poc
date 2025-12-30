package tenantuser

import (
	"context"
	"time"
)

// MockRepository はTenantUserのモックリポジトリ
type MockRepository struct {
	tenantUsers []*TenantUser
}

// NewMockRepository は新しいモックリポジトリを作成する
func NewMockRepository() *MockRepository {
	// モックデータを初期化
	tenantUsers := []*TenantUser{
		// WorkspaceUser: wsu-001 (User 01) の所属テナント
		{
			ID:              "tu-001",
			TenantID:        "tenant-001",
			WorkspaceUserID: "wsu-001",
			Role:            RoleAdmin,
			CreatedAt:       time.Now().Add(-10 * 24 * time.Hour), // 10日前
		},
		{
			ID:              "tu-002",
			TenantID:        "tenant-002",
			WorkspaceUserID: "wsu-001",
			Role:            RoleMember,
			CreatedAt:       time.Now().Add(-8 * 24 * time.Hour), // 8日前
		},
		{
			ID:              "tu-003",
			TenantID:        "tenant-003",
			WorkspaceUserID: "wsu-001",
			Role:            RoleViewer,
			CreatedAt:       time.Now().Add(-5 * 24 * time.Hour), // 5日前
		},
	}

	return &MockRepository{
		tenantUsers: tenantUsers,
	}
}

// FindByWorkspaceUserID はWorkspaceUserIDでTenantUserのリストを取得する
func (r *MockRepository) FindByWorkspaceUserID(ctx context.Context, workspaceUserID string) ([]*TenantUser, error) {
	result := []*TenantUser{}
	for _, tu := range r.tenantUsers {
		if tu.WorkspaceUserID == workspaceUserID {
			result = append(result, tu)
		}
	}
	return result, nil
}
