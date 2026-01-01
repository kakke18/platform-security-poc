import { useAuth } from '~/features/auth';
import { useWorkspaceUsers, WorkspaceUsersList } from '~/features/user';
import { Layout } from '~/components/Layout';

export function UsersPage() {
  const { user: auth0User, loading: authLoading } = useAuth();
  const {
    users,
    loading: usersLoading,
    error: usersError,
    nextPageToken,
    loadMore,
  } = useWorkspaceUsers();

  if (authLoading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!auth0User) {
    return null;
  }

  return (
    <Layout>
      <div className="max-w-4xl">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Workspace Users</h1>

        <WorkspaceUsersList
          users={users}
          loading={usersLoading}
          error={usersError}
          nextPageToken={nextPageToken}
          onLoadMore={loadMore}
        />
      </div>
    </Layout>
  );
}
