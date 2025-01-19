const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

module.exports = {
  entry: "./src/index.ts", // Entry point for TypeScript
  output: {
    filename: "bundle.js", // Output file name
    path: path.resolve(__dirname, "dist"),
    clean: true, // Clean the output folder before each build
  },
  resolve: {
    extensions: [".ts", ".js"], // Resolve TypeScript and JavaScript files
  },
  module: {
    rules: [
      {
        test: /\.ts$/, // Apply ts-loader to TypeScript files
        use: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },
  plugins: [
    new CleanWebpackPlugin(), // Clean dist folder before each build
    new HtmlWebpackPlugin({
      template: "./public/index.html", // Generate an index.html based on the template
    }),
  ],
  devServer: {
    static: {
      directory: path.join(__dirname, "public"), // Serve static files from the public directory
    },
    port: 3000,
    open: true,
    historyApiFallback: true,
    client: {
      logging: "verbose",
    },
  },
  mode: "development",
};
