const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

module.exports = {
	entry: "./src/index.ts",
	output: {
		filename: "bundle.[contenthash].js",
		path: path.resolve(__dirname, "dist"),
		clean: true,
	},
	module: {
		rules: [
			{
				test: /\.ts$/,
				use: "ts-loader",
				exclude: /node_modules/,
			},
			{
				test: /\.css$/,
				use: [MiniCssExtractPlugin.loader, "css-loader"],
			},
		],
	},
	resolve: {
		extensions: [".ts", ".js"],
	},
	plugins: [
		new MiniCssExtractPlugin({
			filename: "styles.[contenthash].css",
		}),
		new HtmlWebpackPlugin({
			template: "./src/index.html",
			minify: {
				removeComments: true,
				useShortDoctype: true,
				collapseWhitespace: true,
				removeRedundantAttributes: true,
				removeScriptTypeAttributes: true,
			},
		}),
	],
	optimization: {
		minimizer: ["...", new CssMinimizerPlugin()],
	},
};
