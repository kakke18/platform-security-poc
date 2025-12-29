import { type NextRequest, NextResponse } from 'next/server';
import { auth0 } from './libs/auth/auth0';

/**
 * Next.jsのプロキシ関数
 * Auth0認証を適用する
 */
export async function proxy(request: NextRequest) {
  // ヘルスチェック用パスへのリクエストの場合は、認証をスキップして200 OKを返す
  if (request.nextUrl.pathname === '/health') {
    return new NextResponse('OK', { status: 200 });
  }

  // Auth0のミドルウェア処理を実行（認証状態の確認など）
  const authRes = await auth0.middleware(request);

  // /auth パスへのリクエストの場合は、Auth0の処理結果をそのまま返す
  // （ログインページなどの認証関連ページへのアクセス）
  if (request.nextUrl.pathname.startsWith('/auth')) {
    return authRes;
  }

  // ユーザセッションを取得
  const session = await auth0.getSession(request);

  // セッションが存在しない（未認証）場合は、ログインページにリダイレクト
  if (!session) {
    return NextResponse.redirect(new URL('/auth/login', request.nextUrl.origin));
  }

  // セッション情報を更新（ユーザ情報に最終更新時間を追加）
  await auth0.updateSession(request, authRes, {
    ...session,
    user: {
      ...session.user,
      // カスタムユーザデータを追加（最終更新時間）
      updatedAt: Date.now(),
    },
  });

  // 更新されたリクエストヘッダーを含む新しいレスポンスを作成
  const resWithCombinedHeaders = NextResponse.next({
    request: {
      headers: request.headers,
    },
  });

  // Auth0レスポンスからヘッダー（特にCookie）を新しいレスポンスに転送
  // これにより、認証状態が維持される
  authRes.headers.forEach((value: string, key: string) => {
    resWithCombinedHeaders.headers.set(key, value);
  });

  // 最終的なレスポンスを返す
  return resWithCombinedHeaders;
}

/**
 * プロキシが適用されるパスを指定する設定
 * 静的ファイルやメタデータファイルには適用されない
 */
export const config = {
  matcher: [
    /*
     * 以下で始まるパス以外のすべてのリクエストパスにマッチ:
     * - _next/static (静的ファイル)
     * - _next/image (画像最適化ファイル)
     * - favicon.ico, sitemap.xml, robots.txt (メタデータファイル)
     */
    '/((?!_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
  ],
};
