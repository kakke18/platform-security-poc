import { readFile } from 'node:fs/promises';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';
import { defineConfig } from 'vitest/config';

const env = {
  ...JSON.parse(await readFile('./config/.env.test.json').then((data) => data.toString())),
};

export default defineConfig({
  plugins: [tsconfigPaths(), react()],
  test: {
    globals: true,
    env,
    projects: [
      {
        extends: 'vitest.config.mts',
        test: {
          name: 'unit',
          environment: 'jsdom',
          include: ['src/**/*.unit.test.{ts,tsx}'],
        },
      },
      {
        extends: 'vitest.config.mts',
        test: {
          name: 'browser',
          include: ['src/**/*.browser.test.{ts,tsx}'],
          browser: {
            enabled: true,
            // @ts-expect-error - vitest browser provider type issue
            provider: 'playwright',
            instances: [{ browser: 'chromium' }],
            headless: true,
            screenshotFailures: false,
          },
        },
      },
    ],
  },
});
