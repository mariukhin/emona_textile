// This config assumes existing usage of the `airbnb` preset, and turns on some
// additional rules, while disabling others. For a complete listing of rules and
// what they do, check out the docs. See:
// - https://github.com/airbnb/javascript/blob/master/packages/eslint-config-airbnb-base/rules/imports.js
// - https://github.com/benmosher/eslint-plugin-import/blob/master/docs/rules

module.exports = {
  plugins: ["import"],
  rules: {
    // Always prefer ES6 `import` unless explicitly disabled due to comment.
    // - https://github.com/benmosher/eslint-plugin-import/blob/master/docs/rules/no-amc.md
    // - https://github.com/benmosher/eslint-plugin-import/blob/master/docs/rules/no-commonjs.md
    "import/no-amd": 2,
    "import/no-commonjs": 2,

    // It is frequently a convention in typescript to prefer named exports.
    // - https://github.com/benmosher/eslint-plugin-import/blob/master/docs/rules/prefer-default-export.md
    "import/prefer-default-export": 0,

    // This rule does not play nicely with index exports and path aliases, so
    // leave it disabled. See:
    // - https://github.com/benmosher/eslint-plugin-import/blob/master/docs/rules/extensions.md
    "import/extensions": 0,
    "import/no-extraneous-dependencies": 0,
  },
};
