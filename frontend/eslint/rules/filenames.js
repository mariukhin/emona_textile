// For a complete listing of rules and what they do, check out the docs. See:
// - https://github.com/selaux/eslint-plugin-filenames

module.exports = {
  plugins: ['filenames'],
  rules: {
    // The filename should match whatever is being exported. See:
    // - https://github.com/selaux/eslint-plugin-filenames
    'filenames/match-exported': [2, null, '\\..*$'],
  },
};
