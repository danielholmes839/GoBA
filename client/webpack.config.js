const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
    watch: true,
    entry: './src/ui/index.js',
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
            },
            {
                test: /\.vue$/,
                loader: 'vue-loader'
            }
        ]
    },

    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        },
        extensions: ['.js', '.ts']
    }, 

    plugins: [
        // make sure to include the plugin!
        new VueLoaderPlugin()
      ]
}