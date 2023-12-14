import { defineUserConfig } from "vuepress";
import theme from "./theme.js";

export default defineUserConfig({
  base: "/",

  lang: "zh-CN",
  title: "Obcsapi Docs",
  description: "Obcsapi Docs | Obcsapi 文档",

  theme,

  // Enable it with pwa
  // shouldPrefetch: false,
});
