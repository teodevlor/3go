// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when a type is not passed). Various options can be set in this
// configuration file.

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'API Documentation',
  tagline: 'Gogogo API - NestJS',

  url: 'https://your-docs-domain.com',
  baseUrl: '/',

  organizationName: 'delitech',
  projectName: 'nest-lam-loi',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  scripts: ['/js/api-docs-config.js'],

  i18n: {
    defaultLocale: 'vi',
    locales: ['vi'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
          editUrl: undefined,
          showLastUpdateTime: true,
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'API Docs',
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'apiSidebar',
            position: 'left',
            label: 'Tài liệu',
          },
        ],
      },
      footer: {
        style: 'dark',
        copyright: `Copyright © ${new Date().getFullYear()} Gogogo API. Built with Docusaurus.`,
      },
      prism: {
        additionalLanguages: ['json', 'bash', 'http'],
      },
    }),
};

module.exports = config;
