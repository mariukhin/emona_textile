// This config extends from the `jsx-a11y/recommended` preset, and uses
// JSX-specific rules from `eslint-plugin-react`. Aside from the existing config
// values, it turns on some additional rules, while disabling others. For a
// complete listing of rules and what they do, check out the docs. See:
// - https://github.com/yannickcr/eslint-plugin-react
// - https://github.com/evcohen/eslint-plugin-jsx-a11y

module.exports = {
  extends: ['plugin:jsx-a11y/recommended'],
  parserOptions: { ecmaFeatures: { jsx: true } },
  plugins: ['jsx-a11y', 'react'],
  rules: {
    'jsx-a11y/anchor-has-content': 0,
    'jsx-a11y/media-has-caption': 0,
    'jsx-a11y/no-onchange': 0,

    // Enforce boolean value consistency. Prefer the format that aligns with the
    // HTML5 spec and requires less typing.
    'react/jsx-boolean-value': 2,

    // Ensure JSX closing tags are consistent, and not all over the place.
    'react/jsx-closing-bracket-location': [2, 'tag-aligned'],

    // Ensure consistency in property styling. Either all single line or all
    // multi-line.
    'react/jsx-first-prop-new-line': [2, 'multiline'],

    // Keep consistent with the normal indentation style.
    'react/jsx-indent-props': [2, 2],

    // Warn if an element that likely requires a key prop – namely, one present
    // in an array literal or an arrow function expression.
    'react/jsx-key': 2,

    'react/jsx-no-bind': [2, { allowArrowFunctions: true }],

    // Help prevent unexpected rendering of comments into the DOM.
    'react/jsx-no-comment-textnodes': 2,

    // Avoid redundantly declaring or mistakenly overriding prop values. This is
    // sometimes done intentionally, but is usually a rare exception.
    'react/jsx-no-duplicate-props': 2,

    // For many base UI elements (eg. Buttons), it is very useful to pass
    // unhandled props to the DOM. For example: miscellaneous aria attributes,
    // eventhandlers, etc., without requiring manual typing and validation
    // in each UI component.
    'react/jsx-props-no-spreading': 0,

    // Forbid target="_blank" attribute without rel="noreferrer" in links, which
    // may expose a security vulnerability.
    'react/jsx-no-target-blank': 2,

    // This rules can help you locate potential ReferenceErrors resulting from
    // misspellings or missing components. Akin to `no-undef`.
    'react/jsx-no-undef': 2,

    // A fragment is redundant if it contains only one child, or if it is the
    // child of a html element, and is not a keyed fragment.
    'react/jsx-no-useless-fragment': 2,

    // Allow one line JSX expressions only when there is a single child
    // element.
    'react/jsx-one-expression-per-line': [2, { allow: 'single-child' }],

    // Ensure consistent component naming that aligns with community standards.
    // Components should be `NamedLikeThis`, `notLikeThis` and `NOTLIKETHIS`.
    'react/jsx-pascal-case': 2,

    // Don't be super-pedantic about requiring sorted properties. There may be
    // logical orderings that are not alphabetical.
    'react/jsx-sort-props': 0,

    // This is necessary for `no-unused-vars` to work. It ensures that when JSX
    // consumes variables they are marked as "used".
    'react/jsx-uses-vars': 2,

    // Components without children can be self-closed to avoid unnecessary extra
    // closing tag. Ensures consistency and saves space.
    'react/self-closing-comp': 2,

    // Enforce multi-line JSX to be enclosed with (). This provides more legible
    // JSX syntax where tag aligment happens on indentation boundaries for the
    // opening and closing tag.
    'react/jsx-wrap-multilines': 2,

    // Don't let people do silly things. :)
    'react/jsx-equals-spacing': [2, 'never'],
  },
};
