package user

import "time"

// User はシステム内のユーザーを表す
type User struct {
	// ID はユーザーID (例: llu_xxxxx)
	ID string

	// Auth0UserID はAuth0のsubject claim (例: auth0|xxxxx, google-oauth2|xxxxx)
	Auth0UserID string

	// WorkspaceID はこのユーザーが所属するワークスペースID
	// ユーザーは1つのワークスペースにのみ所属する (1:N の関係)
	WorkspaceID string

	// IsPrivileged は特権ユーザー（ワークスペース管理者）かどうか
	// - 特権ユーザーはpassword_only認証が強制される（システム固定）
	// - 特権ユーザーはIP制限から除外される
	// - 特権ユーザーはワークスペースレベルのレートリミットから除外される
	IsPrivileged bool

	// IdPConnectionID はAuth0のSSO Connection ID
	// - 特権ユーザーの場合はnullでなければならない（SSO使用不可）
	// - ワークスペースポリシーがsso_onlyの一般ユーザーの場合は必須
	// - ワークスペースポリシーがsso_and_passwordの場合は任意
	IdPConnectionID *string

	// 基本的なプロフィール情報
	Email string
	Name  string

	// タイムスタンプ
	CreatedAt time.Time
	UpdatedAt time.Time
}
