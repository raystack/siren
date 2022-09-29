const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'Siren',
  tagline: 'Universal data observability toolkit',
  url: 'https://odpf.github.io',
  baseUrl: '/siren/',
  onBrokenLinks: 'throw',
  // trailingSlash: true,
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'odpf',
  projectName: 'siren',
  customFields: {
    apiVersion: 'v1beta1',
    defaultHost: 'http://localhost:8080'
  },
  themeConfig: {
    colorMode: {
      defaultMode: 'light',
      respectPrefersColorScheme: true,
      switchConfig: {
        darkIcon: '☾',
        lightIcon: '☀️',
      },
    },
    navbar: {
      title: 'Siren',
      logo: { src: 'img/logo.svg', },
      items: [
        {
          type: 'doc',
          docId: 'introduction',
          position: 'left',
          label: 'Docs',
        },
        { to: '/blog', label: 'Blog', position: 'left' },
        { to: '/help', label: 'Help', position: 'left' },
        {
          href: 'https://bit.ly/2RzPbtn',
          position: 'right',
          className: 'header-slack-link',
        },
        {
          href: 'https://github.com/odpf/siren',
          className: 'navbar-item-github',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      links: [
        {
          title: 'Products',
          items: [
            { label: 'Meteor', href: 'https://github.com/odpf/meteor' },
            { label: 'Firehose', href: 'https://github.com/odpf/firehose' },
            { label: 'Raccoon', href: 'https://github.com/odpf/raccoon' },
            { label: 'Dagger', href: 'https://odpf.github.io/dagger/' },
          ],
        },
        {
          title: 'Resources',
          items: [
            { label: 'Docs', to: '/docs/introduction' },
            { label: 'Blog', to: '/blog', },
            { label: 'Help', to: '/help', },
          ],
        },
        {
          title: 'Community',
          items: [
            { label: 'Slack', href: 'https://bit.ly/2RzPbtn' },
            { label: 'GitHub', href: 'https://github.com/odpf/siren' }
          ],
        },
      ],
      copyright: `Copyright © 2022-${new Date().getFullYear()} ODPF`,
    },
    prism: {
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
    },
    gtag: {
      trackingID: 'G-XXX',
    },
    announcementBar: {
      id: 'star-repo',
      content: '⭐️ If you like Siren, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/odpf/siren">GitHub</a>! ⭐',
      backgroundColor: '#222',
      textColor: '#eee',
      isCloseable: true,
    },
  },

  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/odpf/siren/edit/master/docs/',
        },
        blog: {
          showReadingTime: true,
          editUrl:
            'https://github.com/odpf/siren/edit/master/docs/blog/',
        },
        theme: {
          customCss: [
            require.resolve('./src/css/theme.css'),
            require.resolve('./src/css/custom.css')
          ],
        },
      },
    ],
  ],
};
