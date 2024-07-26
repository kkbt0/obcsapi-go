import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Obcsapi",
  description: "Obsidian 云存储 API",
  base: "/docs/obcsapi/",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Doc", link: "/md/README" },
      { text: "Docs", link: "//www.ftls.xyz/docs/", target: "_blank" },
    ],
    search: {
      provider: "local",
    },
    sidebar: [
      { text: "概述", link: "/md/README" },
      {
        text: "Go 版本",
        items: [
          { text: "1.概述", link: "/md/go-version/1-概述" },
          { text: "2.运行与部署", link: "/md/go-version/2-运行与部署" },
          { text: "3.配置说明", link: "/md/go-version/3-配置说明" },
          { text: "4.功能使用", link: "/md/go-version/4-功能使用" },
          { text: "5.图床说明", link: "/md/go-version/5-图床说明" },
          { text: "6.通用接口", link: "/md/go-version/6-通用接口" },
          { text: "7.前端说明", link: "/md/go-version/7-前端说明" },
          { text: "8.指令模式", link: "/md/go-version/8-指令模式" },
          { text: "9.自定义脚本", link: "/md/go-version/9-自定义脚本" },
          { text: "10.软件联动", link: "/md/go-version/10-软件联动" },
          { text: "97.缓存说明", link: "/md/go-version/97-缓存说明" },
          { text: "98.开发说明", link: "/md/go-version/98-开发说明" },
          { text: "99.更新记录", link: "/md/go-version/99-更新记录" },
        ],
      },
      {
        text: "Go 版本 Swagger(Scalar)",
        link: "/swagger/swagger",
        target: "_blank",
      },
      {
        text: "Python 版本",
        link: "/md/python-version",
      },
      {
        text: "FAQ",
        link: "/md/faq",
      },
    ],

    socialLinks: [
      {
        icon: {
          svg: '<svg xmlns="http://www.w3.org/2000/svg" fill="#000000" width="800px" height="800px" viewBox="0 0 24 24" role="img"><path d="M11.984 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.016 0zm6.09 5.333c.328 0 .593.266.592.593v1.482a.594.594 0 0 1-.593.592H9.777c-.982 0-1.778.796-1.778 1.778v5.63c0 .327.266.592.593.592h5.63c.982 0 1.778-.796 1.778-1.778v-.296a.593.593 0 0 0-.592-.593h-4.15a.592.592 0 0 1-.592-.592v-1.482a.593.593 0 0 1 .593-.592h6.815c.327 0 .593.265.593.592v3.408a4 4 0 0 1-4 4H5.926a.593.593 0 0 1-.593-.593V9.778a4.444 4.444 0 0 1 4.445-4.444h8.296z"/></svg>',
        },
        link: "https://gitee.com/kkbt/obcsapi-go",
      },
      { icon: "github", link: "https://github.com/kkbt0/obcsapi-go" },
    ],
    footer: {
      copyright:
        '版权所有 © 2024 <a href="https://www.ftls.xyz/" target="_blank">恐咖兵糖</a>',
    },
  },
  // TODO
  ignoreDeadLinks: true,
});
