package me

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	gatewayv1 "github.com/kakke18/platform-security-poc/backend/gen/gateway/v1"
	"github.com/kakke18/platform-security-poc/backend/gen/gateway/v1/gatewayv1connect"
	identityv1 "github.com/kakke18/platform-security-poc/backend/gen/identity/v1"
	userv1 "github.com/kakke18/platform-security-poc/backend/gen/user/v1"
)

// Handler はMeServiceの実装
type Handler struct {
	workspaceUserClient identityv1.WorkspaceUserServiceClient
	tenantUserClient    userv1.TenantUserServiceClient
	identityConn        *grpc.ClientConn
	userConn            *grpc.ClientConn
}

// NewHandler は新しいMeハンドラーを作成する
func NewHandler(identityAPIURL, userAPIURL string) (*Handler, error) {
	// Identity APIへのgRPC接続を確立
	identityConn, err := grpc.NewClient(
		identityAPIURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// User APIへのgRPC接続を確立
	userConn, err := grpc.NewClient(
		userAPIURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		identityConn.Close()
		return nil, err
	}

	return &Handler{
		workspaceUserClient: identityv1.NewWorkspaceUserServiceClient(identityConn),
		tenantUserClient:    userv1.NewTenantUserServiceClient(userConn),
		identityConn:        identityConn,
		userConn:            userConn,
	}, nil
}

// Close はgRPC接続をクローズする
func (h *Handler) Close() error {
	if err := h.identityConn.Close(); err != nil {
		slog.Error("Failed to close identity connection", "error", err)
	}
	if err := h.userConn.Close(); err != nil {
		slog.Error("Failed to close user connection", "error", err)
	}
	return nil
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
	identityCtx := metadata.AppendToOutgoingContext(ctx, "x-auth0-user-id", auth0UserID)
	workspaceUserResp, err := h.workspaceUserClient.GetWorkspaceUser(identityCtx, &identityv1.GetWorkspaceUserRequest{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// User APIからTenantUser情報を取得
	userCtx := metadata.AppendToOutgoingContext(ctx, "x-workspace-user-id", workspaceUserResp.WorkspaceUserId)
	tenantUsersResp, err := h.tenantUserClient.GetTenantUsers(userCtx, &userv1.GetTenantUsersRequest{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// TenantUser情報をTenantUserInfo形式に変換
	tenantInfos := make([]*gatewayv1.TenantUserInfo, len(tenantUsersResp.Users))
	for i, tu := range tenantUsersResp.Users {
		tenantInfos[i] = &gatewayv1.TenantUserInfo{
			TenantId:     tu.TenantId,
			TenantUserId: tu.TenantUserId,
			Role:         convertRole(tu.Role),
		}
	}

	// レスポンスを統合
	return connect.NewResponse(&gatewayv1.GetMeResponse{
		WorkspaceId:     workspaceUserResp.WorkspaceId,
		WorkspaceUserId: workspaceUserResp.WorkspaceUserId,
		Email:           workspaceUserResp.Email,
		Name:            workspaceUserResp.Name,
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
