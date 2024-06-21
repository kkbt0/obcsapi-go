import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Obcsapi",
  description: "Obsidian 云存储 API",
  base: "/docs/obcsapi/",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
    ],

    sidebar: [
      { text: '概述', link: '/md/README' },
      {
        text: 'Go 版本',
        items: [
          { text: '1.概述', link: '/md/go-version/1-概述' },
          { text: '2.运行与部署', link: '/md/go-version/2-运行与部署' },
          { text: '3.配置说明', link: '/md/go-version/3-配置说明' },
          { text: '4.功能使用', link: '/md/go-version/4-功能使用' },
          { text: '5.图床说明', link: '/md/go-version/5-图床说明' },
          { text: '6.通用接口', link: '/md/go-version/6-通用接口' },
          { text: '7.前端说明', link: '/md/go-version/7-前端说明' },
          { text: '8.指令模式', link: '/md/go-version/8-指令模式' },
          { text: '9.自定义脚本', link: '/md/go-version/9-自定义脚本' },
          { text: '10.软件联动', link: '/md/go-version/10-软件联动' },
          { text: '97.缓存说明', link: '/md/go-version/97-缓存说明' },
          { text: '98.开发说明', link: '/md/go-version/98-开发说明' },
          { text: '99.更新记录', link: '/md/go-version/99-更新记录' },
        ]
      },
      {
        text: 'Go 版本 Swagger(Scalar)',
        link: '/swagger/swagger'
      },
      {
        text: 'Python 版本',
        link: '/md/python-version'
      },
      {
        text: 'FAQ',
        link: '/md/faq'
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/kkbt0/obcsapi-go' },
    ]
  },
    // TODO
    ignoreDeadLinks: true
})
