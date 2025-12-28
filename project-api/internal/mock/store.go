package mock

import (
	"sync"
	"time"

	"github.com/google/uuid"
	projectv1 "github.com/kakke18/platform-security-poc2/project-api/pkg/gen/project/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Store はインメモリのモックデータストア
type Store struct {
	mu       sync.RWMutex
	projects map[string]*projectv1.Project
	members  map[string][]*projectv1.ProjectMember // key: project_id
}

// NewStore は新しいモックストアを作成し、初期データを投入
func NewStore() *Store {
	s := &Store{
		projects: make(map[string]*projectv1.Project),
		members:  make(map[string][]*projectv1.ProjectMember),
	}

	// サンプルデータを投入
	now := timestamppb.New(time.Now())

	project1 := &projectv1.Project{
		Id:          "proj_1",
		WorkspaceId: "ws_1",
		Name:        "Sample Project 1",
		Description: "This is a sample project for testing",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	project2 := &projectv1.Project{
		Id:          "proj_2",
		WorkspaceId: "ws_1",
		Name:        "Sample Project 2",
		Description: "Another sample project",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.projects[project1.Id] = project1
	s.projects[project2.Id] = project2

	s.members["proj_1"] = []*projectv1.ProjectMember{
		{
			ProjectId: "proj_1",
			UserId:    "user_1",
			Role:      "owner",
			JoinedAt:  now,
		},
		{
			ProjectId: "proj_1",
			UserId:    "user_2",
			Role:      "member",
			JoinedAt:  now,
		},
	}

	s.members["proj_2"] = []*projectv1.ProjectMember{
		{
			ProjectId: "proj_2",
			UserId:    "user_1",
			Role:      "admin",
			JoinedAt:  now,
		},
	}

	return s
}

// CreateProject はプロジェクトを作成
func (s *Store) CreateProject(workspaceID, name, description string) *projectv1.Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := timestamppb.New(time.Now())
	project := &projectv1.Project{
		Id:          "proj_" + uuid.New().String()[:8],
		WorkspaceId: workspaceID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.projects[project.Id] = project
	return project
}

// GetProject はプロジェクトを取得
func (s *Store) GetProject(id string) (*projectv1.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.projects[id]
	return project, ok
}

// ListProjects はワークスペースのプロジェクト一覧を取得
func (s *Store) ListProjects(workspaceID string) []*projectv1.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var projects []*projectv1.Project
	for _, p := range s.projects {
		if p.WorkspaceId == workspaceID {
			projects = append(projects, p)
		}
	}

	return projects
}

// UpdateProject はプロジェクトを更新
func (s *Store) UpdateProject(id, name, description string) (*projectv1.Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	project, ok := s.projects[id]
	if !ok {
		return nil, false
	}

	project.Name = name
	project.Description = description
	project.UpdatedAt = timestamppb.New(time.Now())

	return project, true
}

// DeleteProject はプロジェクトを削除
func (s *Store) DeleteProject(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.projects[id]; !ok {
		return false
	}

	delete(s.projects, id)
	delete(s.members, id)
	return true
}

// AddMember はプロジェクトメンバーを追加
func (s *Store) AddMember(projectID, userID, role string) *projectv1.ProjectMember {
	s.mu.Lock()
	defer s.mu.Unlock()

	member := &projectv1.ProjectMember{
		ProjectId: projectID,
		UserId:    userID,
		Role:      role,
		JoinedAt:  timestamppb.New(time.Now()),
	}

	s.members[projectID] = append(s.members[projectID], member)
	return member
}

// RemoveMember はプロジェクトメンバーを削除
func (s *Store) RemoveMember(projectID, userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	members, ok := s.members[projectID]
	if !ok {
		return false
	}

	for i, m := range members {
		if m.UserId == userID {
			s.members[projectID] = append(members[:i], members[i+1:]...)
			return true
		}
	}

	return false
}

// ListMembers はプロジェクトメンバー一覧を取得
func (s *Store) ListMembers(projectID string) []*projectv1.ProjectMember {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.members[projectID]
}
