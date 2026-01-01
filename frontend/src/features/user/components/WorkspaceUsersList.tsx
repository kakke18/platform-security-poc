import type { WorkspaceUser } from '~/api/generated/gateway/v1/me_pb';

interface WorkspaceUsersListProps {
  users: WorkspaceUser[];
  loading: boolean;
  error: string | null;
  nextPageToken: string;
  onLoadMore: () => void;
}

export function WorkspaceUsersList({
  users,
  loading,
  error,
  nextPageToken,
  onLoadMore,
}: WorkspaceUsersListProps) {
  if (error) {
    return (
      <div className="bg-red-50 border-l-4 border-red-400 p-4">
        <p className="font-semibold text-red-800">Error</p>
        <p className="text-sm text-red-700 mt-1">{error}</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-4">Workspace Users</h2>

      {loading && users.length === 0 ? (
        <div className="text-gray-600">Loading users...</div>
      ) : (
        <>
          <div className="space-y-3">
            {users.map((user) => (
              <div
                key={user.workspaceUserId}
                className="border border-gray-200 rounded-md p-4 hover:bg-gray-50"
              >
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-semibold text-gray-900">{user.name}</p>
                    <p className="text-sm text-gray-600">{user.email}</p>
                    <p className="text-xs text-gray-500 mt-1">ID: {user.workspaceUserId}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {nextPageToken && (
            <button
              type="button"
              onClick={onLoadMore}
              disabled={loading}
              className="mt-4 w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 disabled:bg-gray-400"
            >
              {loading ? 'Loading...' : 'Load More'}
            </button>
          )}

          {users.length === 0 && !loading && (
            <p className="text-gray-500 text-center py-4">No users found</p>
          )}
        </>
      )}
    </div>
  );
}
