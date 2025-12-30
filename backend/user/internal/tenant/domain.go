package tenant

import "time"

// Tenant はプロダクト環境のテナントを表すドメインモデル
type Tenant struct {
	// ID はテナントID (UUID)
	ID string

	// WorkspaceID は所属するワークスペースID
	WorkspaceID string

	// Name はテナント名
	Name string

	// CreatedAt は作成日時
	CreatedAt time.Time
}
