const path = require('path');

const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = (isProd, getFilename) => ({
  plugins: [
    new HTMLWebpackPlugin({
      title: 'Emona Textile',
      template: path.resolve(__dirname, '../public/index.html'),
      minify: {
        collapseWhitespace: isProd,
      },
      favicon: path.resolve(__dirname, '../public/favicon.ico'),
    }),
    new CleanWebpackPlugin(),
    new MiniCssExtractPlugin({
      filename: getFilename('css'),
    }),
  ],
});
