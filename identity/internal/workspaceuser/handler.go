package workspaceuser

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/kakke18/platform-security-poc/identity/gen/identity/v1"
	"github.com/kakke18/platform-security-poc/identity/internal/workspace"
)

// Handler はWorkspaceUserServiceの実装
type Handler struct {
	workspaceUserRepo Repository
	workspaceRepo     workspace.Repository
}

// NewHandler は新しいWorkspaceUserハンドラーを作成する
func NewHandler(workspaceUserRepo Repository, workspaceRepo workspace.Repository) *Handler {
	return &Handler{
		workspaceUserRepo: workspaceUserRepo,
		workspaceRepo:     workspaceRepo,
	}
}

// GetWorkspaceUser は現在のユーザーのWorkspace情報を取得する
func (h *Handler) GetWorkspaceUser(
	ctx context.Context,
	req *connect.Request[identityv1.GetWorkspaceUserRequest],
) (*connect.Response[identityv1.GetWorkspaceUserResponse], error) {
	// ヘッダーからAuth0ユーザーIDを取得（Gatewayで検証済み）
	auth0UserID := req.Header().Get("X-Auth0-User-ID")
	if auth0UserID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// Auth0ユーザーIDでWorkspaceUserを取得
	workspaceUser, err := h.workspaceUserRepo.FindByAuth0UserID(ctx, auth0UserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// WorkspaceIDでWorkspaceを取得
	workspace, err := h.workspaceRepo.FindByID(ctx, workspaceUser.WorkspaceID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&identityv1.GetWorkspaceUserResponse{
		WorkspaceId:     workspace.ID,
		WorkspaceUserId: workspaceUser.ID,
		Email:           workspaceUser.Email,
		Name:            workspaceUser.Name,
	}), nil
}
