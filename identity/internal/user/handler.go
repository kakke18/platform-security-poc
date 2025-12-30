package user

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/kakke18/platform-security-poc/identity/gen/identity/v1"
)

// Handler はUserServiceの実装
type Handler struct {
	repo Repository
}

// NewHandler は新しいユーザーハンドラーを作成する
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetMe は現在のユーザー情報をデータベースから取得する
func (h *Handler) GetMe(
	ctx context.Context,
	req *connect.Request[identityv1.GetMeRequest],
) (*connect.Response[identityv1.GetMeResponse], error) {
	// ヘッダーからユーザーIDを取得（Gatewayで検証済み）
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// Auth0ユーザーIDでデータベースからユーザーを取得
	user, err := h.repo.FindByAuth0UserID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&identityv1.GetMeResponse{
		UserId: user.Auth0UserID,
		Email:  user.Email,
		Name:   user.Name,
	}), nil
}

// UpdateMe は現在のユーザーのプロフィールを更新する
func (h *Handler) UpdateMe(
	ctx context.Context,
	req *connect.Request[identityv1.UpdateMeRequest],
) (*connect.Response[identityv1.UpdateMeResponse], error) {
	// ヘッダーからユーザーIDを取得（Gatewayで検証済み）
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// データベースから現在のユーザーを取得
	user, err := h.repo.FindByAuth0UserID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 提供されたフィールドを更新
	if req.Msg.Name != nil {
		user.Name = *req.Msg.Name
	}

	// データベースに更新を保存
	if err := h.repo.Update(ctx, user); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&identityv1.UpdateMeResponse{
		UserId: user.Auth0UserID,
		Email:  user.Email,
		Name:   user.Name,
	}), nil
}
