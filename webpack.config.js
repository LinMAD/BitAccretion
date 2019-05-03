/* globals __dirname process */
'use strict';

const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    devtool: 'source-map',
    entry: './resources/app.jsx',
    output: {
        path: path.join(__dirname, 'public'),
        publicPath: '',
        filename: 'vizceral.bundle.js'
    },
    resolve: {
        extensions: ['.jsx', '.js'],
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                exclude: /node_modules/,
                loader: 'babel-loader',
            },
            {
                test: /\.html$/,
                use: [ {
                    loader: 'html-loader',
                    options: {
                        minimize: true
                    }
                }],
            },
            {test: /\.woff2?$/, loader: 'url-loader?limit=10000&mimetype=application/font-woff'},
            {test: /\.otf$/, loader: 'file-loader'},
            {test: /\.ttf$/, loader: 'file-loader'},
            {test: /\.eot$/, loader: 'file-loader'},
            {test: /\.svg$/, loader: 'file-loader'},
            {test: /\.css$/, loader: 'style-loader!css-loader'}
        ]
    },
    plugins: [
        new webpack.DefinePlugin({
            'process.env': {
                'NODE_ENV': JSON.stringify('production')
            }
        }),
        new webpack.ProvidePlugin({
            jQuery: 'jquery',
            $: 'jquery'
        }),
        new webpack.DefinePlugin({
            __HIDE_DATA__: !!process.env.HIDE_DATA
        }),
        new HtmlWebpackPlugin({
            title: 'BitAccretion - vizceral',
            template: './resources/index.html',
            favicon: './resources/favicon.ico',
            inject: true
        })
    ]
};
