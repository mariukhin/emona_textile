// node modules
const path = require('path');
// config
const config = require('./prod.config');

const DIST_FOLDER = `dist`;

module.exports = {
  ...config,
  optimization: {
    ...config.optimization,
    splitChunks: {
      ...config.optimization.splitChunks,
      minChunks: 5,
      minSize: 100 * 1000,
      maxSize: 250 * 1000,
      enforceSizeThreshold: 250 * 1000,
    },
  },
  output: {
    ...config.output,
    path: path.resolve(process.cwd(), DIST_FOLDER),
  },
  plugins: [
    ...config.plugins,
  ],
};
