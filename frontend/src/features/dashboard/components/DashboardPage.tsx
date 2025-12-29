import { useAuth } from '~/features/auth';
import { UserInfo, useUser } from '~/features/user';

export function DashboardPage() {
  const { user: auth0User, loading: authLoading, error: authError } = useAuth();
  const { user, loading: userLoading, error: userError } = useUser();

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
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-md p-8">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <a
              href="/auth/logout"
              className="bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600"
            >
              Logout
            </a>
          </div>

          <div className="space-y-6">
            {userLoading && (
              <div className="bg-blue-50 border-l-4 border-blue-400 p-4">
                <p className="text-blue-800">Loading user information...</p>
              </div>
            )}

            {user && <UserInfo user={user} />}

            {(authError || userError) && (
              <div className="bg-red-50 border-l-4 border-red-400 p-4">
                <p className="font-semibold text-red-800">Error</p>
                <p className="text-sm text-red-700 mt-1">{authError || userError}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
