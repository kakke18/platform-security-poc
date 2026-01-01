import { Layout } from '~/components/Layout';
import { PageHeader } from '~/components/PageHeader';
import { useAuth } from '~/features/auth';
import { useWorkspaceUsers, WorkspaceUsersList } from '~/features/user';

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
      <PageHeader title="Workspace Users" />

      <div className="max-w-4xl px-8">
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
