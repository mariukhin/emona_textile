// For complete listing of rules and what they do, check out the docs. See:
// - https://github.com/typescript-eslint/typescript-eslint/tree/master/packages/eslint-plugin

module.exports = {
  overrides: [
    {
      files: ['*.tsx', '*.ts'],
      extends: [
        'plugin:@typescript-eslint/eslint-recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:@typescript-eslint/recommended-requiring-type-checking',
      ],
      rules: {
        // JSDOC-based type comments are unnecessary in typescript files.
        // This rule is only turned off for .ts/.tsx files. See:
        // - https://eslint.org/docs/rules/valid-jsdoc
        'valid-jsdoc': 0,
      },
    },
    {
      files: ['*.stories.tsx', '*.stories.ts'],
      rules: {
        '@typescript-eslint/explicit-module-boundary-types': 0,
      },
    },
    {
      files: ['*.test.tsx', '*.test.ts'],
      rules: {
        // May be needed for accessing TS private members when testing. See:
        // - https://eslint.org/docs/rules/dot-notation
        'dot-notation': 0,
      },
    },
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    createDefaultProgram: true,
    project: './tsconfig.json',
    tsconfigRootDir: './',
  },
  plugins: ['@typescript-eslint', 'import'],
  rules: {
    // Disable things that are checked by Typescript
    // Checked by Typescript - ts(2378)
    'getter-return': 'off',
    // Checked by Typescript - ts(2300)
    'no-dupe-args': 'off',
    // Checked by Typescript - ts(1117)
    'no-dupe-keys': 'off',
    // Checked by Typescript - ts(7027)
    'no-unreachable': 'off',
    // Checked by Typescript - ts(2367)
    'valid-typeof': 'off',
    // Checked by Typescript - ts(2588)
    'no-const-assign': 'off',
    // Checked by Typescript - ts(2588)
    'no-new-symbol': 'off',
    // Checked by Typescript - ts(2376)
    'no-this-before-super': 'off',
    // This is checked by Typescript using the option `strictNullChecks`.
    'no-undef': 'off',
    // This is already checked by Typescript.
    'no-dupe-class-members': 'off',
    // This is already checked by Typescript.
    'no-redeclare': 'off',
    // TS checker handles this
    'import/no-unresolved': 0,
    // TS checker handles this
    'import/named': 0,
    // The spread operator/rest parameters should be prefered in Typescript.
    'prefer-rest-params': 'error',
    'prefer-spread': 'error',

    '@typescript-eslint/interface-name-prefix': 0,
    '@typescript-eslint/no-empty-interface': 0,
    '@typescript-eslint/no-misused-promises': 0,
    '@typescript-eslint/no-non-null-assertion': 0,
    '@typescript-eslint/no-use-before-define': 0,
    '@typescript-eslint/unbound-method': 0,

    // handled by prettier
    '@typescript-eslint/indent': 0,

    // handled by `eslint-plugin-import`
    '@typescript-eslint/no-var-requires': 0,
    '@typescript-eslint/no-require-imports': 0,

    // It's sometimes useful to allow nullish values that would result in an
    // empty string to be interpolated, rather than exhaustively checking
    // beforehand. See:
    // - https://github.com/typescript-eslint/typescript-eslint/blob/master/packages/eslint-plugin/docs/rules/restrict-template-expressions.md
    '@typescript-eslint/restrict-template-expressions': [
      2,
      {
        allowNumber: true,
        allowNullish: true,
      },
    ],
  },
};
