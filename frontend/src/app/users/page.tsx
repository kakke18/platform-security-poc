'use client';

import { withPageAuthRequired } from '@auth0/nextjs-auth0/client';
import { UsersPage } from '~/features/users';

export default withPageAuthRequired(UsersPage, {
  returnTo: '/users',
});
