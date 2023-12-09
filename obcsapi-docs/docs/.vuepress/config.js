import { defineUserConfig } from 'vuepress'
import { defaultTheme } from 'vuepress'

export default defineUserConfig({
  lang: 'zh-CN',
  title: 'Obsidian Cloud Storage API 文档',
  description: 'ObcsapiGo 文档',
  markdown:{
    headers: {
      level: [1,2, 3, 4],
    }
  },
  sidebarDepth: 5,
  theme: defaultTheme({
    sidebar: [
      {
        text: 'Go版本',
        link: '/md/go-version/1-概述.md',
        children: [
          {
            text: '1. 概述',
            link: '/md/go-version/1-概述.md',
          },
          {
            text: '2. 运行与部署',
            link: '/md/go-version/2-运行与部署.md',
          },
          {
           text: '3. 配置说明',
           link: '/md/go-version/3-配置说明.md', 
          },{
            text: '4. 功能使用',
            link: '/md/go-version/4-功能使用.md',
          },{
            text: '5. 图床说明',
            link: '/md/go-version/5-图床说明.md',
          },{
            text: '6. 通用接口',
            link: '/md/go-version/6-通用接口.md',
          },{
            text: '7. 前端说明',
            link: '/md/go-version/7-前端说明.md',
          },{
            text: '8. 指令模式',
            link: '/md/go-version/8-指令模式.md',
          },{
            text: '9. 自定义脚本',
            link: '/md/go-version/9-自定义脚本.md',
          },{
            text: '10. 软件联动',
            link: '/md/go-version/10-软件联动.md',
          },{
            text: '97. 缓存说明',
            link: '/md/go-version/97-缓存说明.md',
          },{
            text: '98. 开发说明',
            link: '/md/go-version/98-开发说明.md',
          },{
            text: '99. 更新记录',
            link: '/md/go-version/99-更新记录.md',
          }
        ]
    },
    {
      text: 'Python版本',
      link: '/md/python-version.md',
    },
    {
      text: 'FAQ',
      link: '/md/faq.md',
    }
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
  }),
})

