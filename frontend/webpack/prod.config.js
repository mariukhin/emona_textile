// config
const base = require('./base.config');

module.exports = {
  ...base,
  mode: 'production',
  optimization: {
    splitChunks: {
      chunks: 'all',
      // minSize: 50 * 1000,
      // minRemainingSize: 0,
      // maxSize: 200 * 1000,
      // minChunks: 1,
      // maxAsyncRequests: 30,
      // maxInitialRequests: 30,
      // automaticNameDelimiter: '~',
      // enforceSizeThreshold: 200 * 1000,
    },
  },
};
