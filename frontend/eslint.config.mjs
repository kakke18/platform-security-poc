import nextPlugin from '@next/eslint-plugin-next';
import typescriptParser from '@typescript-eslint/parser';
import reactPlugin from 'eslint-plugin-react';
import hooksPlugin from 'eslint-plugin-react-hooks';

const config = [
  {
    ignores: ['.next/*', 'node_modules/*', 'pnpm-lock.yaml', 'src/api/generated/**'],
  },
  {
    languageOptions: {
      parser: typescriptParser,
    },
    files: ['**/*.ts', '**/*.tsx'],
    plugins: {
      react: reactPlugin,
      'react-hooks': hooksPlugin,
      '@next/next': nextPlugin,
    },
    rules: {
      ...reactPlugin.configs['jsx-runtime'].rules,
      ...hooksPlugin.configs.recommended.rules,
      ...nextPlugin.configs.recommended.rules,
      ...nextPlugin.configs['core-web-vitals'].rules,
      '@next/next/no-img-element': 'error',
      // eslint-plugin-tailwindcss has not been introduced because it is not yet compatible with Tailwind CSS v4.
      // https://github.com/francoismassart/eslint-plugin-tailwindcss/issues/325
      // ...tailwind.configs["flat/recommended"],
    },
  },
  {
    languageOptions: {
      parser: typescriptParser,
    },
    files: ['**/*.ts', '**/*.tsx'],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              regex: '^swr(/.*)?',
              message: 'Please use ~/libs/swr for swr imports',
            },
          ],
        },
      ],
    },
  },
  {
    languageOptions: {
      parser: typescriptParser,
    },
    files: ['src/libs/swr/*.ts'],
    rules: {
      'no-restricted-imports': ['off'],
    },
  },
  {
    languageOptions: {
      parser: typescriptParser,
    },
    files: ['**/*.ts', '**/*.tsx'],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              regex: '^~/features/(.+)/(.+)',
              message:
                'features配下のファイルは直接importせず、features/*/index.tsからexportしたものを使用してください',
            },
          ],
        },
      ],
    },
  },
];

export default config;
