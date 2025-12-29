import { redirect } from 'next/navigation';
import { auth0 } from '~/libs/auth/auth0';

export default async function Home() {
  const session = await auth0.getSession();

  // Already logged in, redirect to dashboard
  if (session) {
    redirect('/dashboard');
  }

  // Not logged in, redirect to login
  redirect('/auth/login');
}
