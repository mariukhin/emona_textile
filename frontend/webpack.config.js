const path = require('path');
const { merge } = require('webpack-merge');

const devServer = require('./webpack/devServer');
const loaders = require('./webpack/loaders');
const optimization = require('./webpack/optimization');
const plugins = require('./webpack/plugins');
const resolve = require('./webpack/resolve');

const isDev = process.env.NODE_ENV === 'development';
const isProd = !isDev;

const getFilename = ext => (isDev ? `[name].${ext}` : `[name].[hash].${ext}`);

module.exports = merge([
  {
    context: path.resolve(__dirname, 'src'),
    mode: isDev ? 'development' : isProd && 'production',
    bail: isProd,
    entry: {
      main: ['@babel/polyfill', './index.jsx'],
    },
    output: {
      filename: getFilename('js'),
      path: path.resolve(__dirname, 'dist'),
      publicPath: '/',
    },
    devtool: isDev ? 'eval-source-map' : '',
  },
  resolve(),
  optimization(isProd),
  devServer(isDev),
  plugins(isProd, getFilename),
  loaders(isDev),
]);
