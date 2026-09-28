import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  // Main documentation sidebar
  tutorialSidebar: [
    'intro',
    'installation',
    'quickstart',
    'why-fp-go',
    'comparison',
    {
      type: 'category',
      label: 'Migration',
      items: [
        'migration/index',
        'migration/v1-to-v2',
        'migration/interop',
      ],
    },
    {
      type: 'category',
      label: 'Concepts',
      items: [
        'concepts/pure-functions',
        'concepts/monads',
        'concepts/composition',
        'concepts/effects-and-io',
        'concepts/higher-kinded-types',
        'concepts/zen-of-go',
      ],
    },
    'faq',
    'glossary',
    'design-kit',
  ],

  // API Reference sidebar
  apiSidebar: [
    'api-reference',
  ],
};

export default sidebars;

// Made with Bob
