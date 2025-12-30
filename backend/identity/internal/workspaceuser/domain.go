package workspaceuser

import "time"

// WorkspaceUser はWorkspace内のユーザーを表すドメインモデル
type WorkspaceUser struct {
	// ID はワークスペースユーザーID (UUID)
	ID string

	// WorkspaceID は所属するワークスペースID
	WorkspaceID string

	// Auth0UserID はAuth0のユーザーID (sub claim)
	Auth0UserID string

	// Email はメールアドレス
	Email string

	// Name は表示名
	Name string

	// CreatedAt は作成日時
	CreatedAt time.Time
}
