// This config extends from the `react/recommended` preset, and turns on some
// additional rules, while disabling others. For a complete listing of rules and
// what they do, check out the docs. See:
// - https://github.com/facebook/react/tree/master/packages/eslint-plugin-react-hooks
// - https://github.com/yannickcr/eslint-plugin-react

module.exports = {
  extends: ['plugin:react/recommended'],
  parserOptions: { ecmaFeatures: { jsx: true } },
  plugins: ['react', 'react-hooks'],
  rules: {
    'react/destructuring-assignment': [
      2,
      'always',
      { ignoreClassFields: true },
    ],

    // For ES6 class components and stateless components this isn't necessary.
    'react/display-name': 0,

    'react/jsx-filename-extension': [2, { extensions: ['.jsx', '.tsx'] }],

    // It's frequently more convenient (and less typing!) to treat children as
    // a prop value.
    'react/no-children-prop': 2,

    // Updating the state after a component mounts will trigger a second
    // render() call and can lead to property/layout thrashing.
    'react/no-did-mount-set-state': 2,

    // Updating the state after a component updates will trigger a second
    // render() call and can lead to property/layout thrashing.
    'react/no-did-update-set-state': 2,

    // In general, one component per file allows for strong isolation, easier
    // testing and easier stubbing/mocking. This strategy is alright for full
    // class components but becomes annoying when dealing with stateless
    // components or indeed any small utility functions you have for generating
    // some set components you need (e.g. an iterator for a map call). Overall
    // this winds up being more a hinderance than a help.
    'react/no-multi-comp': 0,

    // Help prevent common typos.
    'react/no-typos': 2,

    // Stateless components are generally preferable because they're easier
    // to reason about, simpler to write and (eventually) more performant.
    'react/prefer-stateless-function': [1, { ignorePureComponents: true }],

    // This is turned off in favor of using typescript to annotate stateless
    // function components.
    'react/prop-types': 0,

    // In functional components, default prop values may be specified as ES6
    // default parameters. See:
    // - https://github.com/yannickcr/eslint-plugin-react/blob/master/docs/rules/require-default-props.md
    'react/require-default-props': [2, { ignoreFunctionalComponents: true }],

    // Let people order properties however they want.
    'react/sort-comp': 0,

    // These enforce the Rules of Hooks.
    // See: https://reactjs.org/docs/hooks-rules.html
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',
  },
  settings: { react: { version: 'detect' } },
};
