import type { GetMeResponse } from '~/api/generated/identity/v1/user_pb';

interface UserInfoProps {
  user: GetMeResponse;
}

export function UserInfo({ user }: UserInfoProps) {
  return (
    <div>
      <h2 className="text-xl font-semibold text-gray-900 mb-3">
        User Information (from Backend API)
      </h2>
      <div className="bg-blue-50 border border-blue-200 p-4 rounded-md">
        <div className="space-y-2 text-gray-800">
          <p>
            <span className="font-semibold">User ID:</span> {user.userId}
          </p>
          <p>
            <span className="font-semibold">Email:</span> {user.email}
          </p>
          <p>
            <span className="font-semibold">Name:</span> {user.name}
          </p>
        </div>
      </div>
      <p className="text-sm text-gray-600 mt-2">
        ✓ User information verified by Identity API via JWT authentication
      </p>
    </div>
  );
}
