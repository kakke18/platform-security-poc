import { useEffect, useState } from 'react';
import { meServiceClient } from '~/api/connect-client';
import type { WorkspaceUser } from '~/api/generated/gateway/v1/me_pb';

export interface UseWorkspaceUsersResult {
  users: WorkspaceUser[];
  loading: boolean;
  error: string | null;
  nextPageToken: string;
  loadMore: () => Promise<void>;
}

export function useWorkspaceUsers(): UseWorkspaceUsersResult {
  const [users, setUsers] = useState<WorkspaceUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [nextPageToken, setNextPageToken] = useState('');

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setLoading(true);
        const response = await meServiceClient.listWorkspaceUsers({
          pageSize: 10,
          pageToken: '',
        });

        setUsers(response.users);
        setNextPageToken(response.nextPageToken);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch workspace users');
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, []);

  const loadMore = async () => {
    if (nextPageToken && !loading) {
      try {
        setLoading(true);
        const response = await meServiceClient.listWorkspaceUsers({
          pageSize: 10,
          pageToken: nextPageToken,
        });

        setUsers((prev) => [...prev, ...response.users]);
        setNextPageToken(response.nextPageToken);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch workspace users');
      } finally {
        setLoading(false);
      }
    }
  };

  return { users, loading, error, nextPageToken, loadMore };
}
