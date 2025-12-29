import { useUser } from '@auth0/nextjs-auth0/client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

interface UseAuthReturn {
  user: {
    name?: string;
    email?: string;
    sub?: string;
  } | null;
  loading: boolean;
  error: string | null;
}

export function useAuth(): UseAuthReturn {
  const router = useRouter();
  const { user, error: auth0Error, isLoading } = useUser();

  useEffect(() => {
    // Auth0のローディングが完了するまで待つ
    if (isLoading) {
      return;
    }

    // ユーザーがいない場合はログインにリダイレクト
    if (!user) {
      router.push('/auth/login');
    }
  }, [user, isLoading, router]);

  return {
    user: user || null,
    loading: isLoading,
    error: auth0Error?.message || null,
  };
}
