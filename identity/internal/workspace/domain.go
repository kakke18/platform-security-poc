package workspace

import "time"

// Workspace はテナントを跨ぐ組織を表すドメインモデル
type Workspace struct {
	// ID はワークスペースID (UUID)
	ID string

	// Name はワークスペース名
	Name string

	// CreatedAt は作成日時
	CreatedAt time.Time
}
