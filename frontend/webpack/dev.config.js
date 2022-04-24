// node modules
const path = require('path');
// plugins
const { HotModuleReplacementPlugin } = require('webpack');
// config
const base = require('./base.config');

module.exports = {
  ...base,
  cache: true,
  devtool: 'inline-source-map',
  devServer: {
    contentBase: path.resolve(process.cwd(), 'src'),
    hot: true,
    port: 3000,
    publicPath: '/',
    historyApiFallback: true,
    compress: true,
  },
  optimization: {
    removeAvailableModules: false,
    removeEmptyChunks: false,
    splitChunks: false,
  },
  plugins: [
    ...base.plugins,
    new HotModuleReplacementPlugin(),
  ],
};
