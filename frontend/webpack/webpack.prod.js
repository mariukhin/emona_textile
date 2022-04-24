// node modules
const path = require('path');
const fs = require('fs');
const parser = require('xml2json');
// config
const config = require('@upstox/ui-config-react/webpack/prod.config');
// plugins
const { DefinePlugin } = require('webpack');
const CopyPlugin = require('copy-webpack-plugin');

// consts
const ENV = process.env.ENV || 'UAT';
const LOCAL = process.env.LOCAL || 'FALSE';
const DIST_FOLDER = `dist/${ENV}`;
const POM_FILE = fs.readFileSync(path.resolve(__dirname, '../pom.xml'));
const VERSION = parser.toJson(POM_FILE, { object: true }).project.version;

console.info(`ENV=${ENV}, VERSION=${VERSION}, LOCAL=${LOCAL}`);

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
  devtool: ENV === 'PROD' ? false : 'source-map',
  ...(ENV === 'PROD' &&
    LOCAL === 'FALSE' && {
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
    }),
  plugins: [
    ...config.plugins,
    new CopyPlugin({
      patterns: [
        {
          from: path.resolve(process.cwd(), 'static'),
          to: path.resolve(process.cwd(), DIST_FOLDER),
        },
      ],
    }),
    new DefinePlugin({
      'process.env.LOCAL': JSON.stringify(LOCAL),
      'process.env.ENV': JSON.stringify(ENV),
      'process.env.VERSION': JSON.stringify(VERSION),
    }),
  ],
};
