import { useEffect, useState } from 'react';
import { meServiceClient } from '~/api/connect-client';
import type { GetMeResponse } from '~/api/generated/gateway/v1/me_pb';

interface UseUserReturn {
  user: GetMeResponse | null;
  loading: boolean;
  error: string | null;
}

export function useUser(): UseUserReturn {
  const [user, setUser] = useState<GetMeResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchUser() {
      setLoading(true);
      setError(null);
      try {
        // Connect clientは自動的にAuth0トークンを付与する
        const response = await meServiceClient.getMe({});
        setUser(response);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch user info');
        console.error('Failed to fetch user info:', err);
      } finally {
        setLoading(false);
      }
    }

    fetchUser();
  }, []);

  return { user, loading, error };
}
