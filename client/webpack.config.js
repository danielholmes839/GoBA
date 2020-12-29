module.exports = {
    watch: true,
    entry: './src/index.ts',
    output: {
        path: __dirname + '/dist',
        filename: 'bundle.js',
        publicPath: '/dist/'
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                exclude: '/node_modules/',
                use: {
                    loader: 'ts-loader'
                }
            }
        ]
    },

    resolve: {
        extensions: ['.ts']
    }
}