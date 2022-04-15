const path = require('path');

module.exports = isDev => ({
  devServer: isDev
    ? {
        port: 3000,
        contentBase: path.resolve(__dirname, 'dist'),
        hot: true,
        overlay: true,
        historyApiFallback: true,
        compress: true,
      }
    : {},
});
