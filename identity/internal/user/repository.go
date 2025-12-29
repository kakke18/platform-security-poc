package user

import "context"

// Repository はユーザーデータアクセスのインターフェース
type Repository interface {
	// FindByAuth0UserID はAuth0ユーザーIDでユーザーを取得する
	FindByAuth0UserID(ctx context.Context, auth0UserID string) (*User, error)

	// Update はユーザー情報を更新する
	Update(ctx context.Context, user *User) error
}
