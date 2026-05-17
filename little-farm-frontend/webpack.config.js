const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
	stats: "minimal",
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
		],
	},
	resolve: {
		extensions: [".ts", ".js"],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: "./src/index.html",
			minify: {
				html5: true,
				useShortDoctype: true,
				collapseWhitespace: true,
				minifyJS: true,
				minifyCSS: true,
				minifyURLs: false,
				removeComments: true,
				removeOptionalTags: true,
				removeAttributeQuotes: true,
				removeEmptyAttributes: true,
				removeRedundantAttributes: true,
				removeScriptTypeAttributes: true,
				removeStyleLinkTypeAttributes: true,
			},
		}),
	],
	optimization: {
		minimizer: ["..."],
	},
};
