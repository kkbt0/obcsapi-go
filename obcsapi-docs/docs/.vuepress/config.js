module.exports = {
  lang: 'zh-CN',
  title: 'Obsidian Cloud Storage API 文档',
  description: 'ObcsapiGo 文档',
  host: '0.0.0.0',
  port: 1313,
  markdown:{
    headers: {
      level: [1,2, 3, 4],
    }
  },
  sidebarDepth: 5,
  themeConfig: {
    sidebar: [
      {
        title: 'Go版本',   // 必要的
        collapsable: false, // 可选的, 默认值是 true,
        sidebarDepth: 10,    // 可选的, 默认值是 1
        children: [
          ['/md/go-version/1-概述.md', '1. 概述'],
          ['/md/go-version/2-运行与部署.md','2. 运行与部署'],
          ['/md/go-version/3-配置说明.md','3. 配置说明'],
          ['/md/go-version/4-功能使用.md','4. 功能使用'],
          ['/md/go-version/5-图床说明.md','5. 图床说明'],
          ['/md/go-version/6-通用接口.md','6. 通用接口'],
          ['/md/go-version/7-前端说明.md','7. 前端说明'],
          ['/md/go-version/8-指令模式.md','8. 指令模式'],
          ['/md/go-version/9-自定义脚本.md','9. 自定义脚本'],
          ['/md/go-version/10-软件联动.md','10. 软件联动'],
          ['/md/go-version/97-缓存说明.md','97. 缓存说明'],
          ['/md/go-version/98-开发说明.md','98. 开发说明'],
          ['/md/go-version/99-更新记录.md','99. 更新记录']
        ]
      },
  ],
  navbar: [
    {
        text: '首页',
        link: '/',
    },
    {
        text: 'Go 版本源码',
        link: 'https://gitee.com/kkbt/obcsapi-go',
    }
]
  },
}

