package tenant

import (
	"context"
	"fmt"
	"time"
)

// MockRepository はTenantのモックリポジトリ
type MockRepository struct {
	tenants map[string]*Tenant
}

// NewMockRepository は新しいモックリポジトリを作成する
func NewMockRepository() *MockRepository {
	// モックデータを初期化
	tenants := map[string]*Tenant{
		"tenant-001": {
			ID:          "tenant-001",
			WorkspaceID: "ws-001",
			Name:        "Production",
			CreatedAt:   time.Now().Add(-20 * 24 * time.Hour), // 20日前
		},
		"tenant-002": {
			ID:          "tenant-002",
			WorkspaceID: "ws-001",
			Name:        "Staging",
			CreatedAt:   time.Now().Add(-15 * 24 * time.Hour), // 15日前
		},
		"tenant-003": {
			ID:          "tenant-003",
			WorkspaceID: "ws-001",
			Name:        "Development",
			CreatedAt:   time.Now().Add(-10 * 24 * time.Hour), // 10日前
		},
	}

	return &MockRepository{
		tenants: tenants,
	}
}

// FindByID はIDでTenantを取得する
func (r *MockRepository) FindByID(ctx context.Context, id string) (*Tenant, error) {
	tenant, ok := r.tenants[id]
	if !ok {
		return nil, fmt.Errorf("tenant not found: %s", id)
	}
	return tenant, nil
}
