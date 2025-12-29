import { Auth0Client } from '@auth0/nextjs-auth0/server';

const getRequiredEnv = (key: string): string => {
  const value = process.env[key];
  if (!value) {
    throw new Error(`Missing required environment variable: ${key}`);
  }
  return value;
};

export const auth0 = new Auth0Client({
  secret: getRequiredEnv('AUTH0_SECRET'),
  appBaseUrl: getRequiredEnv('AUTH0_BASE_URL'),
  clientId: getRequiredEnv('AUTH0_CLIENT_ID'),
  clientSecret: getRequiredEnv('AUTH0_CLIENT_SECRET'),
  authorizationParameters: {
    audience: process.env.AUTH0_AUDIENCE,
  },
});
