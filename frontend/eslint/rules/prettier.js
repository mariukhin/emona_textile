// This config extends from the `prettier/@typescript-eslint` preset, and  turns
// on some additional rules, while disabling others. For a complete listing of
// rules and what they do, check out the docs. See:
// - https://github.com/prettier/eslint-plugin-prettier

module.exports = {
  extends: ['prettier/@typescript-eslint'],
  plugins: ['prettier'],
  rules: { 'prettier/prettier': 'error' },
  overrides: [
    {
      files: ['*.js', '*.jsx', '*.ts', '*.tsx'],
      rules: {
        // The Airbnb preset enforces this rule, and while it is preferred, Prettier
        // may enforce a linebreak in some contexts for nicer formatting, so we
        // relax the rules here.
        // - https://eslint.org/docs/rules/implicit-arrow-linebreak
        'implicit-arrow-linebreak': 0,

        'function-paren-newline': ['error', 'multiline-arguments'],

        // Allow prettier to manage this.
        // - https://eslint.org/docs/rules/indent
        indent: 0,

        // Let prettier handle max line length. See:
        // - https://eslint.org/docs/rules/max-len
        'max-len': 0,

        // This style is enforced by Pretter, and it clases with the value of the
        // AirBnB preset for this rule. Since this style is not configurable in
        // Prettier, we relax the strictness of the rule here. See:
        // - https://eslint.org/docs/rules/object-curly-newline
        'object-curly-newline': [
          'error',
          {
            ExportDeclaration: { consistent: true, multiline: true },
            ImportDeclaration: { consistent: true, multiline: true },
            ObjectExpression: { consistent: true, multiline: true },
            ObjectPattern: { consistent: true, multiline: true },
          },
        ],

        // This is the style enforced by Pretter, which is essentially the opposite
        // value of the AirBnB preset for this rule. Since this style is not
        // configurable in Prettier, we reset the value here. See:
        // - https://eslint.org/docs/rules/operator-linebreak
        // - https://github.com/prettier/prettier/issues/3806
        'operator-linebreak': [
          'error',
          'after',
          { overrides: { ':': 'before', '?': 'before' } },
        ],
      },
    },
  ],
};
