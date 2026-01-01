'use client';

import { withPageAuthRequired } from '@auth0/nextjs-auth0/client';
import { MePage } from '~/features/me';

export default withPageAuthRequired(MePage, {
  returnTo: '/me',
});
