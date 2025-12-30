package tenantuser

import "time"

// Role はテナント内でのユーザーのロール
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

// TenantUser はTenant内のユーザーを表すドメインモデル
type TenantUser struct {
	// ID はテナントユーザーID (UUID)
	ID string

	// TenantID は所属するテナントID
	TenantID string

	// WorkspaceUserID は対応するワークスペースユーザーID
	WorkspaceUserID string

	// Role はテナント内でのロール
	Role Role

	// CreatedAt は作成日時
	CreatedAt time.Time
}
