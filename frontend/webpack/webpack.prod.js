// node modules
const path = require('path');
// config
const config = require('./prod.config');
// plugins
const CopyPlugin = require('copy-webpack-plugin');

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
  devtool: {
      module: {
        ...config.module,
        rules: [
          ...config.module.rules,
          {
            test: /\.tsx?$/,
            exclude: [/node_modules/],
            use: [
              {
                loader: 'babel-loader',
                options: {
                  plugins: ['@babel/plugin-proposal-class-properties'],
                },
              },
              {
                loader: 'ts-loader',
                options: {
                  transpileOnly: true,
                  configFile: path.resolve(process.cwd(), 'tsconfig.json'),
                },
              },
            ],
          },
        ],
      },
    },
  plugins: [
    ...config.plugins,
    new CopyPlugin({
      patterns: [
        {
          from: path.resolve(process.cwd(), 'static'),
          to: path.resolve(process.cwd(), DIST_FOLDER),
        },
      ],
    })
  ],
};
