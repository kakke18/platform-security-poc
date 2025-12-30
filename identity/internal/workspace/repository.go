package workspace

import "context"

// Repository はWorkspaceのリポジトリインターフェース
type Repository interface {
	// FindByID はIDでWorkspaceを取得する
	FindByID(ctx context.Context, id string) (*Workspace, error)
}
