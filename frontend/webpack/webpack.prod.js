// node modules
const path = require('path');
// config
const config = require('./prod.config');

const DIST_FOLDER = `dist`;

module.exports = {
  ...config,
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
