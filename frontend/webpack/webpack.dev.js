// node modules
const path = require('path');
// config
const config = require('./dev.config');
// plugins
// const { DefinePlugin } = require('webpack');

const DIST_FOLDER = `dist`;

module.exports = {
  ...config,
  devServer: {
    ...config.devServer,
    host: 'dev.emona.com',
    port: 443,
    open: true,
  },
  output: {
    ...config.output,
    path: path.resolve(process.cwd(), DIST_FOLDER),
  },
  plugins: [
    ...config.plugins,
  ],
};
