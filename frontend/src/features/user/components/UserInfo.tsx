import { useState } from 'react';
import { type GetMeResponse, Role } from '~/api/generated/gateway/v1/me_pb';

interface UserInfoProps {
  user: GetMeResponse;
  accessToken: string | null;
}

function getRoleLabel(role: Role): string {
  switch (role) {
    case Role.ADMIN:
      return 'Admin';
    case Role.MEMBER:
      return 'Member';
    case Role.VIEWER:
      return 'Viewer';
    default:
      return 'Unknown';
  }
}

export function UserInfo({ user, accessToken }: UserInfoProps) {
  const [copied, setCopied] = useState(false);

  const handleCopyToken = () => {
    if (accessToken) {
      navigator.clipboard.writeText(accessToken);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div>
      <h2 className="text-xl font-semibold text-gray-900 mb-3">
        User Information (from Backend API)
      </h2>
      <div className="bg-blue-50 border border-blue-200 p-4 rounded-md">
        <div className="space-y-2 text-gray-800">
          <p>
            <span className="font-semibold">Workspace ID:</span> {user.workspaceId}
          </p>
          <p>
            <span className="font-semibold">Workspace User ID:</span> {user.workspaceUserId}
          </p>
          <p>
            <span className="font-semibold">Email:</span> {user.email}
          </p>
          <p>
            <span className="font-semibold">Name:</span> {user.name}
          </p>
        </div>
      </div>

      {/* Access Token Section */}
      <div className="mt-4 bg-yellow-50 border border-yellow-200 p-4 rounded-md">
        <div className="flex justify-between items-center mb-2">
          <h3 className="text-lg font-semibold text-gray-900">Access Token</h3>
          {accessToken && (
            <button
              type="button"
              onClick={handleCopyToken}
              className="bg-blue-500 text-white py-1 px-3 rounded-md hover:bg-blue-600 text-sm"
            >
              {copied ? 'Copied!' : 'Copy Token'}
            </button>
          )}
        </div>
        {accessToken ? (
          <div className="bg-gray-900 text-green-400 p-3 rounded font-mono text-xs overflow-x-auto">
            {accessToken}
          </div>
        ) : (
          <p className="text-red-600 text-sm">No access token available</p>
        )}
        <p className="text-xs text-gray-600 mt-2">
          Use this token with:{' '}
          <code className="bg-gray-200 px-1 rounded">
            AUTH0_TOKEN=&quot;...&quot; ./scripts/test-gateway.sh
          </code>
        </p>
      </div>

      {user.tenants && user.tenants.length > 0 && (
        <div className="mt-4 bg-green-50 border border-green-200 p-4 rounded-md">
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Tenants</h3>
          <div className="space-y-3">
            {user.tenants.map((tenant) => (
              <div key={tenant.tenantId} className="bg-white p-3 rounded border border-gray-200">
                <p className="text-gray-800">
                  <span className="font-semibold">Tenant ID:</span> {tenant.tenantId}
                </p>
                <p className="text-gray-800">
                  <span className="font-semibold">Tenant User ID:</span> {tenant.tenantUserId}
                </p>
                <p className="text-gray-800">
                  <span className="font-semibold">Role:</span>{' '}
                  <span className="inline-block px-2 py-1 text-xs font-semibold rounded bg-blue-100 text-blue-800">
                    {getRoleLabel(tenant.role)}
                  </span>
                </p>
              </div>
            ))}
          </div>
        </div>
      )}

      <p className="text-sm text-gray-600 mt-2">
        ✓ User information aggregated from Identity API and User Service via JWT authentication
      </p>
    </div>
  );
}
