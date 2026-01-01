package me

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	gatewayv1 "github.com/kakke18/platform-security-poc/backend/gen/gateway/v1"
	"github.com/kakke18/platform-security-poc/backend/gen/gateway/v1/gatewayv1connect"
	identityv1 "github.com/kakke18/platform-security-poc/backend/gen/identity/v1"
	"github.com/kakke18/platform-security-poc/backend/gen/identity/v1/identityv1connect"
	userv1 "github.com/kakke18/platform-security-poc/backend/gen/user/v1"
	"github.com/kakke18/platform-security-poc/backend/gen/user/v1/userv1connect"
)

// Handler はMeServiceの実装
type Handler struct {
	workspaceUserClient identityv1connect.WorkspaceUserServiceClient
	tenantUserClient    userv1connect.TenantUserServiceClient
}

// NewHandler は新しいMeハンドラーを作成する
// gRPCプロトコル（HTTP/2 over cleartext）を使用してバックエンドサービスと通信
func NewHandler(identityAPIURL, userAPIURL string) *Handler {
	// HTTP/2クライアントを作成（h2c: HTTP/2 Cleartext）
	h2cClient := &http.Client{
		Transport: &http2.Transport{
			// h2c (HTTP/2 without TLS) を許可
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				// TLSダイヤルの代わりに通常のダイヤルを使用
				return net.Dial(network, addr)
			},
		},
	}

	return &Handler{
		workspaceUserClient: identityv1connect.NewWorkspaceUserServiceClient(
			h2cClient,
			identityAPIURL,
			connect.WithGRPC(), // gRPCプロトコルを使用
		),
		tenantUserClient: userv1connect.NewTenantUserServiceClient(
			h2cClient,
			userAPIURL,
			connect.WithGRPC(), // gRPCプロトコルを使用
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

// ListWorkspaceUsers はワークスペース内のユーザー一覧を取得する
func (h *Handler) ListWorkspaceUsers(
	ctx context.Context,
	req *connect.Request[gatewayv1.ListWorkspaceUsersRequest],
) (*connect.Response[gatewayv1.ListWorkspaceUsersResponse], error) {
	// JWTミドルウェアで設定されたAuth0ユーザーIDを取得
	auth0UserID := req.Header().Get("X-Auth0-User-ID")
	if auth0UserID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	// Identity APIからワークスペース内のユーザー一覧を取得
	listReq := connect.NewRequest(&identityv1.ListWorkspaceUsersRequest{
		PageSize:  req.Msg.PageSize,
		PageToken: req.Msg.PageToken,
	})
	listReq.Header().Set("X-Auth0-User-ID", auth0UserID)

	listResp, err := h.workspaceUserClient.ListWorkspaceUsers(ctx, listReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Identity APIのWorkspaceUserをGatewayのWorkspaceUserに変換
	users := make([]*gatewayv1.WorkspaceUser, len(listResp.Msg.Users))
	for i, u := range listResp.Msg.Users {
		users[i] = &gatewayv1.WorkspaceUser{
			WorkspaceUserId: u.WorkspaceUserId,
			Email:           u.Email,
			Name:            u.Name,
		}
	}

	return connect.NewResponse(&gatewayv1.ListWorkspaceUsersResponse{
		Users:         users,
		NextPageToken: listResp.Msg.NextPageToken,
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
