import { Layout } from '~/components/Layout';
import { PageHeader } from '~/components/PageHeader';
import { useAuth } from '~/features/auth';
import { UserInfo, useUser } from '~/features/user';

export function MePage() {
  const { user: auth0User, loading: authLoading } = useAuth();
  const { user, loading: userLoading, error: userError, accessToken } = useUser();

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
      <PageHeader title="My Profile" />

      <div className="max-w-4xl px-8">
        <div className="bg-white rounded-lg shadow-md p-8">
          <div className="space-y-6">
            {userLoading && (
              <div className="bg-blue-50 border-l-4 border-blue-400 p-4">
                <p className="text-blue-800">Loading user information...</p>
              </div>
            )}

            {user && <UserInfo user={user} accessToken={accessToken} />}

            {userError && (
              <div className="bg-red-50 border-l-4 border-red-400 p-4">
                <p className="font-semibold text-red-800">Error</p>
                <p className="text-sm text-red-700 mt-1">{userError}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
}
