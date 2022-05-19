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
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        exclude: [/node_modules/],
        use: {
          loader: 'ts-loader',
          options: {
            transpileOnly: true,
            configFile: path.resolve(process.cwd(), 'tsconfig.json'),
          },
        },
      },
      {
        test: /\.(cur|png|jpg|jpeg|svg|woff|woff2|eot|ttf)$/i,
        use: [
          {
            loader: 'file-loader',
            options: {
              name: '[name]~[hash].[ext]',
              outputPath: 'assets',
            },
          },
        ],
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
    ],
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
