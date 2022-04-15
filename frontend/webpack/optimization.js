const OptimizeCssAssetsWebpackPlugin = require('optimize-css-assets-webpack-plugin');
const TerserWebpackPlugin = require('terser-webpack-plugin');

module.exports = isProd => {
  const optimization = {
    splitChunks: {
      chunks: 'all',
    },
  };

  if (isProd) {
    optimization.minimizer = [
      new OptimizeCssAssetsWebpackPlugin(),
      new TerserWebpackPlugin(),
    ];
  }

  return { optimization };
};
