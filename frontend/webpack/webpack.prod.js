// node modules
const path = require('path');
// config
const config = require('./prod.config');
const CopyPlugin = require('copy-webpack-plugin');

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
        test: /\.html$/,
        use: ['html-loader'],
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
    ],
  },
  output: {
    ...config.output,
    path: path.resolve(process.cwd(), DIST_FOLDER),
  },
  plugins: [
    ...config.plugins,
    new CopyPlugin({
      patterns: [
        {
          from: path.resolve(process.cwd(), 'src/assets'),
          to: path.resolve(process.cwd(), `${DIST_FOLDER}/assets`),
        },
        {
          from: path.resolve(process.cwd(), 'src/fonts'),
          to: path.resolve(process.cwd(), `${DIST_FOLDER}/fonts`),
        },
      ],
    }),
  ],
};
