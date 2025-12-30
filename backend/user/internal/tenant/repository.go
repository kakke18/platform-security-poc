package tenant

import "context"

// Repository はTenantのリポジトリインターフェース
type Repository interface {
	// FindByID はIDでTenantを取得する
	FindByID(ctx context.Context, id string) (*Tenant, error)
}
