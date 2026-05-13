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
    {
      type: 'category',
      label: 'v2 Core Types',
      items: [
        {
          type: 'category',
          label: 'Essential Types',
          items: [
            'v2/either',
            'v2/io',
            'v2/ioresult',
            'v2/ioeither',
            'v2/iooption',
            'v2/result',
            'v2/effect',
          ],
        },
        {
          type: 'category',
          label: 'Reader Types',
          items: [
            'v2/reader',
            'v2/readereither',
            'v2/readerio',
            'v2/readerioeither',
            'v2/readerioresult',
            'v2/readeroption',
          ],
        },
        {
          type: 'category',
          label: 'State & Advanced',
          items: [
            'v2/state',
            'v2/statereaderioeither',
            'v2/lazy',
            'v2/constant',
            'v2/identity',
            'v2/endomorphism',
          ],
        },
      ],
    },
    {
      type: 'category',
      label: 'Collections',
      items: [
        {
          type: 'category',
          label: 'Array',
          items: [
            'v2/collections/array',
            'v2/collections/array-ap',
            'v2/collections/array-eq',
            'v2/collections/array-find',
            'v2/collections/array-monoid',
            'v2/collections/array-sort',
            'v2/collections/array-uniq',
            'v2/collections/array-zip',
            'v2/collections/nonempty-array',
          ],
        },
        {
          type: 'category',
          label: 'Record',
          items: [
            'v2/collections/record',
            'v2/collections/record-ap',
            'v2/collections/record-chain',
            'v2/collections/record-conversion',
            'v2/collections/record-eq',
            'v2/collections/record-monoid',
            'v2/collections/record-ord',
            'v2/collections/record-traverse',
          ],
        },
        {
          type: 'category',
          label: 'Sequence & Traverse',
          items: [
            'v2/collections/sequence-traverse',
          ],
        },
      ],
    },
    {
      type: 'category',
      label: 'Utilities',
      items: [
        'v2/utilities/function',
        'v2/utilities/pipe-flow',
        'v2/utilities/compose',
        'v2/utilities/bind-curry',
        'v2/utilities/predicate',
        'v2/utilities/boolean',
        'v2/utilities/number',
        'v2/utilities/string',
        'v2/utilities/tuple',
        'v2/utilities/eq',
        'v2/utilities/ord',
        'v2/utilities/semigroup',
        'v2/utilities/monoid',
        'v2/utilities/magma',
      ],
    },
    {
      type: 'category',
      label: 'Advanced',
      items: [
        'advanced/patterns',
        'advanced/type-theory',
        'advanced/performance',
        'advanced/architecture',
      ],
    },
  ],

  // Recipes sidebar
  recipesSidebar: [
    'recipes-index',
    {
      type: 'category',
      label: 'Error Handling',
      items: [
        'recipes/validation',
        'recipes/error-recovery',
        'recipes/error-handling',
        'recipes/retry',
      ],
    },
    {
      type: 'category',
      label: 'Data Processing',
      items: [
        'recipes/data-transformation',
        'recipes/filtering-mapping',
        'recipes/aggregation',
        'recipes/parsing',
      ],
    },
    {
      type: 'category',
      label: 'I/O Operations',
      items: [
        'recipes/file-operations',
        'recipes/http-requests',
        'recipes/parallel-tasks',
      ],
    },
    {
      type: 'category',
      label: 'Composition Patterns',
      items: [
        'recipes/dependency-injection',
        'recipes/pipelines',
        'recipes/middleware',
      ],
    },
    {
      type: 'category',
      label: 'Testing',
      items: [
        'recipes/testing-pure',
        'recipes/testing-effects',
      ],
    },
  ],
};

export default sidebars;

// Made with Bob
