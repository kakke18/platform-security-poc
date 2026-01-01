package workspaceuser

import "context"

// Repository はWorkspaceUserのリポジトリインターフェース
type Repository interface {
	// FindByAuth0UserID はAuth0ユーザーIDでWorkspaceUserを取得する
	FindByAuth0UserID(ctx context.Context, auth0UserID string) (*WorkspaceUser, error)

	// ListByWorkspaceID はワークスペースIDでWorkspaceUserの一覧を取得する
	ListByWorkspaceID(ctx context.Context, workspaceID string, pageSize int32, pageToken string) ([]*WorkspaceUser, string, error)
}
