package me

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	gatewayv1 "github.com/kakke18/platform-security-poc/backend/gateway/gen/gateway/v1"
	"github.com/kakke18/platform-security-poc/backend/gateway/gen/gateway/v1/gatewayv1connect"
	identityv1 "github.com/kakke18/platform-security-poc/backend/identity/gen/identity/v1"
	"github.com/kakke18/platform-security-poc/backend/identity/gen/identity/v1/identityv1connect"
	userv1 "github.com/kakke18/platform-security-poc/backend/user/gen/user/v1"
	"github.com/kakke18/platform-security-poc/backend/user/gen/user/v1/userv1connect"
)

// Handler はMeServiceの実装
type Handler struct {
	workspaceUserClient identityv1connect.WorkspaceUserServiceClient
	tenantUserClient    userv1connect.TenantUserServiceClient
}

// NewHandler は新しいMeハンドラーを作成する
func NewHandler(identityAPIURL, userAPIURL string) *Handler {
	return &Handler{
		workspaceUserClient: identityv1connect.NewWorkspaceUserServiceClient(
			http.DefaultClient,
			identityAPIURL,
		),
		tenantUserClient: userv1connect.NewTenantUserServiceClient(
			http.DefaultClient,
			userAPIURL,
		),
	}
}

// Ensure Handler implements gatewayv1connect.MeServiceHandler
var _ gatewayv1connect.MeServiceHandler = (*Handler)(nil)

// GetMe は現在認証されているユーザーの全情報を取得する
func (h *Handler) GetMe(
	ctx context.Context,
	req *connect.Request[gatewayv1.GetMeRequest],
) (*connect.Response[gatewayv1.GetMeResponse], error) {
	// JWTミドルウェアで設定されたAuth0ユーザーIDを取得
	auth0UserID := req.Header().Get("X-Auth0-User-ID")
	if auth0UserID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// Identity APIからWorkspaceUser情報を取得
	workspaceUserReq := connect.NewRequest(&identityv1.GetWorkspaceUserRequest{})
	workspaceUserReq.Header().Set("X-Auth0-User-ID", auth0UserID)

	workspaceUserResp, err := h.workspaceUserClient.GetWorkspaceUser(ctx, workspaceUserReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// User APIからTenantUser情報を取得
	tenantUsersReq := connect.NewRequest(&userv1.GetTenantUsersRequest{})
	tenantUsersReq.Header().Set("X-Workspace-User-ID", workspaceUserResp.Msg.WorkspaceUserId)

	tenantUsersResp, err := h.tenantUserClient.GetTenantUsers(ctx, tenantUsersReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// TenantUser情報をTenantUserInfo形式に変換
	tenantInfos := make([]*gatewayv1.TenantUserInfo, len(tenantUsersResp.Msg.Users))
	for i, tu := range tenantUsersResp.Msg.Users {
		tenantInfos[i] = &gatewayv1.TenantUserInfo{
			TenantId:     tu.TenantId,
			TenantUserId: tu.TenantUserId,
			Role:         convertRole(tu.Role),
		}
	}

	// レスポンスを統合
	return connect.NewResponse(&gatewayv1.GetMeResponse{
		WorkspaceId:     workspaceUserResp.Msg.WorkspaceId,
		WorkspaceUserId: workspaceUserResp.Msg.WorkspaceUserId,
		Email:           workspaceUserResp.Msg.Email,
		Name:            workspaceUserResp.Msg.Name,
		Tenants:         tenantInfos,
	}), nil
}

// convertRole はUser ServiceのRoleをGatewayのRoleに変換する
func convertRole(role userv1.Role) gatewayv1.Role {
	switch role {
	case userv1.Role_ROLE_ADMIN:
		return gatewayv1.Role_ROLE_ADMIN
	case userv1.Role_ROLE_MEMBER:
		return gatewayv1.Role_ROLE_MEMBER
	case userv1.Role_ROLE_VIEWER:
		return gatewayv1.Role_ROLE_VIEWER
	default:
		return gatewayv1.Role_ROLE_UNSPECIFIED
	}
}
