package workspaceuser

import "context"

// Repository はWorkspaceUserのリポジトリインターフェース
type Repository interface {
	// FindByAuth0UserID はAuth0ユーザーIDでWorkspaceUserを取得する
	FindByAuth0UserID(ctx context.Context, auth0UserID string) (*WorkspaceUser, error)
}
