package tenantuser

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/kakke18/platform-security-poc/backend/user/gen/user/v1"
)

// Handler はTenantUserServiceの実装
type Handler struct {
	repo Repository
}

// NewHandler は新しいTenantUserハンドラーを作成する
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetTenantUsers は現在のWorkspaceUserが所属するTenantUserのリストを取得する
func (h *Handler) GetTenantUsers(
	ctx context.Context,
	req *connect.Request[userv1.GetTenantUsersRequest],
) (*connect.Response[userv1.GetTenantUsersResponse], error) {
	// ヘッダーからWorkspaceUserIDを取得（Gatewayで検証済み）
	workspaceUserID := req.Header().Get("X-Workspace-User-ID")
	if workspaceUserID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// WorkspaceUserIDでTenantUserのリストを取得
	tenantUsers, err := h.repo.FindByWorkspaceUserID(ctx, workspaceUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// ドメインモデルをProtoメッセージに変換
	protoUsers := make([]*userv1.TenantUser, len(tenantUsers))
	for i, tu := range tenantUsers {
		protoUsers[i] = &userv1.TenantUser{
			TenantId:     tu.TenantID,
			TenantUserId: tu.ID,
			Role:         roleToProto(tu.Role),
		}
	}

	return connect.NewResponse(&userv1.GetTenantUsersResponse{
		Users: protoUsers,
	}), nil
}

// roleToProto はドメインモデルのRoleをProtoのRoleに変換する
func roleToProto(role Role) userv1.Role {
	switch role {
	case RoleAdmin:
		return userv1.Role_ROLE_ADMIN
	case RoleMember:
		return userv1.Role_ROLE_MEMBER
	case RoleViewer:
		return userv1.Role_ROLE_VIEWER
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}
