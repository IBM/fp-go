import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'fp-go',
  tagline: 'Functional programming for Go',
  favicon: 'img/fp-go-logo.png',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: 'https://ibm.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/fp-go/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'IBM', // Usually your GitHub org/user name.
  projectName: 'fp-go', // Usually your repo name.

  onBrokenLinks: 'warn',
  // Broken-anchor checker reads markdown heading slugs from the MDX AST and
  // does not see JSX-supplied id attributes on <Section> headings. The
  // anchors work at runtime, so silence the false positives.
  onBrokenAnchors: 'ignore',

  // codapi-js — turns <codapi-snippet> tags in MDX into live editors.
  // Points at the local codapi server during dev; swap the URL for a hosted
  // instance before deploying.
  stylesheets: [
    {href: 'https://unpkg.com/@antonz/codapi@0.20.0/dist/snippet.css'},
  ],
  scripts: [
    {src: 'https://unpkg.com/@antonz/codapi@0.20.0/dist/snippet.js', defer: true},
  ],

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          lastVersion: 'current',
          versions: {
            current: {
              label: 'v2.2.82 (latest)',
              path: '',
              banner: 'none',
            },
            '1.0.0': {
              label: 'v1.x (legacy)',
              path: '1.0.0',
            },
          },
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    // Algolia DocSearch configuration
    algolia: {
      // The application ID provided by Algolia
      appId: 'YOUR_APP_ID',
      
      // Public API key: it is safe to commit it
      apiKey: 'YOUR_SEARCH_API_KEY',
      
      // Index name
      indexName: 'fp-go',
      
      // Optional: see doc section below
      contextualSearch: true,
      
      // Optional: Specify domains where the navigation should occur through window.location instead on history.push
      // Useful when our Algolia config crawls multiple documentation sites and we want to navigate with window.location.href to them.
      // externalUrlRegex: 'external\\.com|domain\\.com',
      
      // Optional: Replace parts of the item URLs from Algolia. Useful when using the same search index for multiple deployments using a different baseUrl.
      // replaceSearchResultPathname: {
      //   from: '/docs/', // or as RegExp: /\/docs\//
      //   to: '/',
      // },
      
      // Optional: Algolia search parameters
      searchParameters: {},
      
      // Optional: path for search page that enabled by default (`false` to disable it)
      searchPagePath: 'search',
      
      // Optional: whether the insights feature is enabled or not on Docsearch (`false` by default)
      insights: false,
    },
    navbar: {
      title: 'fp-go',
      logo: {
        alt: 'fp-go Logo',
        src: 'img/fp-go-logo.png',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'docSidebar',
          sidebarId: 'apiSidebar',
          position: 'left',
          label: 'API Reference',
        },
        {
          type: 'docSidebar',
          sidebarId: 'recipesSidebar',
          position: 'left',
          label: 'Recipes',
        },
        {
          type: 'docsVersionDropdown',
          position: 'right',
          dropdownItemsAfter: [],
          dropdownActiveClassDisabled: true,
        },
        {
          type: 'search',
          position: 'right',
        },
        {
          href: 'https://github.com/IBM/fp-go',
          label: 'GitHub',
          position: 'right',
          className: 'header-github-link',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Documentation',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/intro',
            },
            {
              label: 'Quick Start',
              to: '/docs/quickstart',
            },
            {
              label: 'Installation',
              to: '/docs/installation',
            },
            {
              label: 'Migration Guide',
              to: '/docs/migration',
            },
          ],
        },
        {
          title: 'Learn',
          items: [
            {
              label: 'Core Concepts',
              to: '/docs/concepts',
            },
            {
              label: 'Recipes',
              to: '/docs/recipes',
            },
            {
              label: 'API Reference',
              to: '/docs/v2/result',
            },
            {
              label: 'FAQ',
              to: '/docs/faq',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/IBM/fp-go',
            },
            {
              label: 'Issues',
              href: 'https://github.com/IBM/fp-go/issues',
            },
            {
              label: 'Discussions',
              href: 'https://github.com/IBM/fp-go/discussions',
            },
            {
              label: 'Contributing',
              href: 'https://github.com/IBM/fp-go/blob/main/CONTRIBUTING.md',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'IBM Open Source',
              href: 'https://ibm.github.io/',
            },
            {
              label: 'Go Package',
              href: 'https://pkg.go.dev/github.com/IBM/fp-go/v2',
            },
            {
              label: 'License',
              href: 'https://github.com/IBM/fp-go/blob/main/LICENSE',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} IBM Corporation. Licensed under Apache License 2.0. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,

};

export default config;
