package service

import (
	"context"

	"connectrpc.com/connect"
	projectv1 "github.com/kakke18/platform-security-poc2/project-api/pkg/gen/project/v1"
	"github.com/kakke18/platform-security-poc2/project-api/pkg/gen/project/v1/projectv1connect"
	"github.com/kakke18/platform-security-poc2/project-api/internal/mock"
)

// ProjectService は ProjectService の実装
type ProjectService struct {
	store *mock.Store
}

// NewProjectService は新しい ProjectService を作成
func NewProjectService(store *mock.Store) projectv1connect.ProjectServiceHandler {
	return &ProjectService{
		store: store,
	}
}

// CreateProject はプロジェクトを作成
func (s *ProjectService) CreateProject(
	ctx context.Context,
	req *connect.Request[projectv1.CreateProjectRequest],
) (*connect.Response[projectv1.CreateProjectResponse], error) {
	project := s.store.CreateProject(
		req.Msg.WorkspaceId,
		req.Msg.Name,
		req.Msg.Description,
	)

	return connect.NewResponse(&projectv1.CreateProjectResponse{
		Project: project,
	}), nil
}

// GetProject はプロジェクトを取得
func (s *ProjectService) GetProject(
	ctx context.Context,
	req *connect.Request[projectv1.GetProjectRequest],
) (*connect.Response[projectv1.GetProjectResponse], error) {
	project, ok := s.store.GetProject(req.Msg.Id)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}

	return connect.NewResponse(&projectv1.GetProjectResponse{
		Project: project,
	}), nil
}

// ListProjects はプロジェクト一覧を取得
func (s *ProjectService) ListProjects(
	ctx context.Context,
	req *connect.Request[projectv1.ListProjectsRequest],
) (*connect.Response[projectv1.ListProjectsResponse], error) {
	projects := s.store.ListProjects(req.Msg.WorkspaceId)

	return connect.NewResponse(&projectv1.ListProjectsResponse{
		Projects: projects,
	}), nil
}

// UpdateProject はプロジェクトを更新
func (s *ProjectService) UpdateProject(
	ctx context.Context,
	req *connect.Request[projectv1.UpdateProjectRequest],
) (*connect.Response[projectv1.UpdateProjectResponse], error) {
	project, ok := s.store.UpdateProject(
		req.Msg.Id,
		req.Msg.Name,
		req.Msg.Description,
	)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}

	return connect.NewResponse(&projectv1.UpdateProjectResponse{
		Project: project,
	}), nil
}

// DeleteProject はプロジェクトを削除
func (s *ProjectService) DeleteProject(
	ctx context.Context,
	req *connect.Request[projectv1.DeleteProjectRequest],
) (*connect.Response[projectv1.DeleteProjectResponse], error) {
	success := s.store.DeleteProject(req.Msg.Id)

	return connect.NewResponse(&projectv1.DeleteProjectResponse{
		Success: success,
	}), nil
}

// AddMember はプロジェクトメンバーを追加
func (s *ProjectService) AddMember(
	ctx context.Context,
	req *connect.Request[projectv1.AddMemberRequest],
) (*connect.Response[projectv1.AddMemberResponse], error) {
	member := s.store.AddMember(
		req.Msg.ProjectId,
		req.Msg.UserId,
		req.Msg.Role,
	)

	return connect.NewResponse(&projectv1.AddMemberResponse{
		Member: member,
	}), nil
}

// RemoveMember はプロジェクトメンバーを削除
func (s *ProjectService) RemoveMember(
	ctx context.Context,
	req *connect.Request[projectv1.RemoveMemberRequest],
) (*connect.Response[projectv1.RemoveMemberResponse], error) {
	success := s.store.RemoveMember(
		req.Msg.ProjectId,
		req.Msg.UserId,
	)

	return connect.NewResponse(&projectv1.RemoveMemberResponse{
		Success: success,
	}), nil
}

// ListMembers はプロジェクトメンバー一覧を取得
func (s *ProjectService) ListMembers(
	ctx context.Context,
	req *connect.Request[projectv1.ListMembersRequest],
) (*connect.Response[projectv1.ListMembersResponse], error) {
	members := s.store.ListMembers(req.Msg.ProjectId)

	return connect.NewResponse(&projectv1.ListMembersResponse{
		Members: members,
	}), nil
}
