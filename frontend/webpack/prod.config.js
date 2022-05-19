// config
const base = require('./base.config');

module.exports = {
  ...base,
  mode: 'production',
  optimization: {
    splitChunks: {
      chunks: 'all',
    },
  },
};
