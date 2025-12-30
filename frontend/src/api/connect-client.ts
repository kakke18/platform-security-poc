import { getAccessToken } from '@auth0/nextjs-auth0/client';
import { createClient, type Interceptor } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { MeService } from './generated/gateway/v1/me_pb';
import { UserService } from './generated/identity/v1/user_pb';

/**
 * Auth0トークンを自動的に付与するインターセプター
 * サーバー側ではauth0.getSession()、クライアント側ではgetAccessToken()を使用
 */
const authInterceptor: Interceptor = (next) => async (req) => {
  let token: string | undefined;

  if (typeof window === 'undefined') {
    // サーバーサイド: セッションからトークンを取得
    const { auth0 } = await import('~/libs/auth/auth0');
    const session = await auth0.getSession();
    token = session?.accessToken as string | undefined;
  } else {
    // クライアントサイド: Auth0 SDKからトークンを取得
    try {
      token = await getAccessToken();
    } catch (error) {
      console.error('Failed to get access token:', error);
    }
  }

  if (token) {
    req.header.set('Authorization', `Bearer ${token}`);
  }

  return await next(req);
};

/**
 * Connect Transport の設定
 * - ベースURL: 環境変数またはデフォルト値
 * - インターセプター: 認証トークンを自動付与
 */
const transport = createConnectTransport({
  baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  interceptors: [authInterceptor],
});

/**
 * UserService Connect クライアント
 * 自動的にAuth0トークンを付与してバックエンドAPIにリクエストを送信
 */
export const userServiceClient = createClient(UserService, transport);

/**
 * MeService Connect クライアント
 * Identity と User Service の情報を統合して返すGateway API
 * 自動的にAuth0トークンを付与してバックエンドAPIにリクエストを送信
 */
export const meServiceClient = createClient(MeService, transport);
