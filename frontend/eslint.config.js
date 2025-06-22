//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config'

export default [
  ...tanstackConfig,
  {
    ignores: [
      'vite.config.js',
      'eslint.config.js',
      'prettier.config.js',
    ],
  },
];
