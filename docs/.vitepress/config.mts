import { defineConfig } from "vitepress";

const socialLinksConfig = [
  {
    icon: "github",
    link: "https://github.com/patppuccin/snipraw",
    ariaLabel: "GitHub",
  },
  {
    icon: "docker",
    link: "https://hub.docker.com/r/patppuccin/snipraw",
    ariaLabel: "Docker",
  },
];

const navBarConfig = [
  { text: "Docs", link: "/docs/" },
  { text: "About", link: "/about" },
  {
    text: "Extras",
    items: [
      { text: "License", link: "/license" },
      { text: "Roadmap", link: "/roadmap" },
      { text: "Changelog", link: "/changelog" },
    ],
  },
];

const docsSidebarConfig = [
  {
    text: "Guide",
    link: "/docs/guide/",
    items: [
      { text: "Installation", link: "/docs/guide/installation" },
      { text: "Command Line", link: "/docs/guide/cli" },
      { text: "Configuration", link: "/docs/guide/configuration" },
      { text: "Keyboard Shortcuts", link: "/docs/guide/keyboard-shortcuts" },
    ],
  },
  {
    text: "Deploy",
    link: "/docs/deploy/",
    items: [
      { text: "Container", link: "/docs/deploy/container" },
      { text: "System Service", link: "/docs/deploy/service" },
      { text: "Reverse Proxy", link: "/docs/deploy/reverse-proxy" },
    ],
  },
];

// https://vitepress.dev/reference/site-config
export default defineConfig({
  srcDir: "content",

  title: "Snipraw",
  description: "A minimal personal code snippet server",
  themeConfig: {
    logo: "/logo.svg",
    siteTitle: "Patppuccin",
    search: { provider: "local" },
    outline: [2, 3],
    nav: navBarConfig,
    sidebar: { "/docs/": docsSidebarConfig },
    socialLinks: socialLinksConfig,
    footer: {
      message: "Released under the <a href='/license'>Apache 2.0 License</a>",
      copyright: `Copyright © ${new Date().getFullYear()} <a href='https://patrickambrose.com'>Patrick Ambrose</a>.`,
    },
  },
  markdown: {
    theme: {
      light: "vitesse-light",
      dark: "vitesse-dark",
    },
  },
});
