
module.exports = {

    publicPath: process.env.NODE_ENV === 'production' ? '/' : '/',

    // outputDir: 在npm run build时 生成文件的目录 type:string, default:'dist'

    outputDir: 'dist/static',

    // pages:{ type:Object,Default:undfind } 
    /*
      构建多页面模式的应用程序.每个“页面”都应该有一个相应的JavaScript条目文件。该值应该是一
      个对象，其中键是条目的名称，而该值要么是指定其条目、模板和文件名的对象，要么是指定其条目
      的字符串，
      注意：请保证pages里配置的路径和文件名 在你的文档目录都存在 否则启动服务会报错的
    */
    // pages: {
    // index: {
    // entry for the page
    // entry: 'src/index/main.js',
    // the source template
    // template: 'public/index.html',
    // output as dist/index.html
    // filename: 'index.html'
    // },
    // when using the entry-only string format,
    // template is inferred to be 'public/subpage.html'
    // and falls back to 'public/index.html' if not found.
    // Output filename is inferred to be 'subpage.html'.
    // subpage: 'src/subpage/main.js'
    // },

    //   lintOnSave：{ type:Boolean default:true } 问你是否使用eslint
    lintOnSave: false,
    // productionSourceMap：{ type:Bollean,default:true } 生产源映射
    // 如果您不需要生产时的源映射，那么将此设置为false可以加速生产构建
    productionSourceMap: false,
    // devServer:{type:Object} 3个属性host,port,https
    // 它支持webPack-dev-server的所有选项

    devServer: {
        port: 8077, // 端口号
        host: '0.0.0.0',
        https: false, // https:{type:Boolean}
        // open: true, //配置自动启动浏览器
        // proxy: 'http://localhost:4000' // 配置跨域处理,只有一个代理
        proxy: {
            '/api': {
                target: '<url>',
                ws: true,
                changeOrigin: true
            },
            '/foo': {
                target: '<other_url>'
            }
        },  // 配置多个代理
    },
}