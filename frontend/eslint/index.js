module.exports = {
  extends: [
    './rules/base',
    './rules/filenames',
    './rules/import',
    './rules/jsx',
    './rules/prettier',
    './rules/react',
    './rules/typescript',
  ],
  env: {
    browser: true,
    node: true,
    es2017: true,
    es6: true,
  },
  settings: {
    'import/resolver': {
      node: {
        paths: ['src'],
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    },
  },
  ignorePatterns: ['**/node_modules/*', '**typedoc-theme/*'],
  rules: {
    '@typescript-eslint/explicit-function-return-type': 0,
    '@typescript-eslint/camelcase': 0,
  },
};
