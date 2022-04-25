// This config extends from the `airbnb` and `eslint:recommended` presets, and
// turns on some additional rules, while disabling others. For a complete
// listing of rules and what they do, check out the docs. See:
// - https://github.com/facebook/react/tree/master/packages/eslint-plugin-react-hooks
// - https://eslint.org/docs/rules/

module.exports = {
  env: { es6: true },
  extends: ['airbnb', 'eslint:recommended'],
  parserOptions: { sourceType: 'module' },
  rules: {
    // Omit parentheses when they have exactly one parameter.
    // - https://eslint.org/docs/rules/arrow-parens
    'arrow-parens': [2, 'as-needed', { requireForBlockBody: true }],

    // Ensure consistent linebreaks across OS's.
    // - https://eslint.org/docs/rules/linebreak-style
    'linebreak-style': [0, 'error', 'windows'],

    // Allow immediately adjacent class members when the expression is one line.
    // - https://eslint.org/docs/rules/lines-between-class-members
    'lines-between-class-members': [
      'error',
      'always',
      { exceptAfterSingleLine: true },
    ],

    // Allow underscores in names at the developers discretion. See:
    // - https://eslint.org/docs/rules/no-underscore-dangle
    'no-underscore-dangle': 0,

    // Enforce good documentation. Reduce the mental burden for when other
    // people try to understand or consume your code. Parameter and return
    // typing are not required, as this is handled via Typescript annotations.
    // See:
    // - http://eslint.org/docs/rules/valid-jsdoc
    'valid-jsdoc': [
      2,
      {
        prefer: { return: 'returns' },
      },
    ],
  },
};
