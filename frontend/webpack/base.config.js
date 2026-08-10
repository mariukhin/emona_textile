// node modules
const path = require('path');
const dotenv = require('dotenv');
// plugins
const { DefinePlugin, ProvidePlugin } = require('webpack');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');

// Load frontend/.env. Variables already present in the environment (CI, shell)
// take precedence, dotenv does not overwrite them.
dotenv.config({ path: path.resolve(process.cwd(), '.env') });

// Keys inlined into the browser bundle. Everything listed here is readable by
// anyone who opens the built JS, so it may only hold public identifiers —
// never a real secret. Restrict these keys in the Google Cloud and EmailJS
// dashboards instead.
const CLIENT_ENV_KEYS = [
  'GOOGLE_MAPS_API_KEY',
  'EMAILJS_SERVICE_ID',
  'EMAILJS_TEMPLATE_ID',
  'EMAILJS_PUBLIC_KEY',
  'GA_TRACKING_ID',
];

const clientEnv = CLIENT_ENV_KEYS.reduce(
  (acc, key) => ({
    ...acc,
    [`process.env.${key}`]: JSON.stringify(process.env[key] || ''),
  }),
  { 'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'PROD') },
);

module.exports = {
  target: 'web',
  mode: 'development',
  cache: true,
  entry: path.resolve(process.cwd(), 'src/index.ts'),
  output: {
    path: path.resolve(process.cwd(), 'dist'),
    publicPath: '/',
    filename: '[name]~[hash].js',
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
              name: '[name].[ext]',
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
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    plugins: [
      new TsconfigPathsPlugin({
        configFile: path.resolve(process.cwd(), 'tsconfig.json'),
      }),
    ],
  },
  plugins: [
    new CleanWebpackPlugin(),
    new HtmlWebpackPlugin({
      template: path.resolve(process.cwd(), 'public/index.html'),
      favicon: path.resolve(process.cwd(), 'public/favicon.ico'),
    }),
    new ForkTsCheckerWebpackPlugin(),
    new DefinePlugin(clientEnv),
    // new ProvidePlugin({
    //   process: 'process/browser',
    // }),
  ],
};
