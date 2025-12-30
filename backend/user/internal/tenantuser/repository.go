package tenantuser

import "context"

// Repository はTenantUserのリポジトリインターフェース
type Repository interface {
	// FindByWorkspaceUserID はWorkspaceUserIDでTenantUserのリストを取得する
	FindByWorkspaceUserID(ctx context.Context, workspaceUserID string) ([]*TenantUser, error)
}
