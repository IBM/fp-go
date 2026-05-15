# fp-go Documentation Site

> Comprehensive documentation for fp-go - Functional Programming for Go

This directory contains the complete Docusaurus-based documentation site for the fp-go library. The site provides guides, API references, recipes, and interactive examples to help developers learn and use functional programming patterns in Go.

## 🎮 Interactive Playground

The documentation includes an interactive code playground powered by [fp-go-sandbox](https://github.com/teyouale/fp-go-sandbox). This allows users to:

- Write and execute fp-go code directly in the browser
- Experiment with examples without local setup
- Test code snippets from the documentation
- Learn by doing with immediate feedback

The sandbox integration enables live code execution for all examples throughout the documentation.

## 🚀 Quick Start

### Prerequisites

- **Node.js**: >= 20.0
- **npm**: >= 8.0

### Installation

```bash
cd fp-go-docs
npm install
```

### Development

Start the development server with hot reload:

```bash
npm start
```

The site will open at [http://localhost:3000](http://localhost:3000).

### Build

Generate static content for production:

```bash
npm run build
```

The static files will be generated in the `build/` directory.

### Serve Production Build

Test the production build locally:

```bash
npm run serve
```

## 📁 Project Structure

```
fp-go-docs/
├── docs/                          # Documentation content
│   ├── intro.md                   # Getting started
│   ├── quickstart.md              # Quick start guide
│   ├── installation.md            # Installation instructions
│   ├── concepts/                  # Core FP concepts
│   │   ├── monads.md
│   │   ├── composition.md
│   │   ├── pure-functions.md
│   │   └── higher-kinded-types.md
│   ├── recipes/                   # Practical examples
│   │   ├── error-handling.md
│   │   ├── validation.md
│   │   ├── data-transformation.md
│   │   └── ...
│   ├── v2/                        # v2 API reference
│   │   ├── result.md
│   │   ├── either.md
│   │   ├── option.md
│   │   └── collections/
│   ├── migration/                 # Migration guides
│   │   ├── v1-to-v2.md
│   │   └── interop.md
│   └── advanced/                  # Advanced topics
│       ├── architecture.md
│       ├── patterns.md
│       └── performance.md
├── versioned_docs/                # Versioned documentation
│   └── version-1.0.0/            # v1.x legacy docs
├── src/                           # React components & theme
│   ├── components/
│   │   ├── content/              # Custom MDX components
│   │   │   ├── PageHeader.tsx
│   │   │   ├── CodeCard.tsx
│   │   │   ├── Callout.tsx
│   │   │   ├── Bench.tsx
│   │   │   └── ...
│   │   ├── HomepageFeatures/
│   │   └── Playground/           # Interactive code playground
│   ├── css/
│   │   └── custom.css            # Custom styling
│   ├── pages/
│   │   ├── index.tsx             # Homepage
│   │   └── playground.tsx        # Playground page
│   └── theme/                     # Theme customizations
├── static/                        # Static assets
│   └── img/
├── docusaurus.config.ts          # Docusaurus configuration
├── sidebars.ts                   # Sidebar structure
├── package.json                  # Dependencies & scripts
└── tsconfig.json                 # TypeScript config
```

## 📝 Writing Documentation

### Creating a New Page

1. Create a new `.md` or `.mdx` file in the appropriate directory under `docs/`
2. Add frontmatter:

```markdown
---
id: my-page
title: My Page Title
sidebar_label: Short Label
description: Brief description for SEO
---

# My Page Title

Content goes here...
```

3. Update `sidebars.ts` to include your new page

### Using Custom Components

The documentation includes custom MDX components for rich content:

#### PageHeader

```mdx
<PageHeader
  eyebrow="Guide · Section 01 / 03"
  title="Error"
  titleAccent="handling."
  lede="Learn how to handle errors functionally using Result and Either types."
  meta={[
    {label: '// Version', value: 'v2.2.82'},
    {label: '// Reading time', value: '5 min'},
  ]}
/>
```

#### CodeCard

```mdx
<CodeCard file="example.go" status="tested">
{`func Example() Result[int] {
    return R.Of(42)
}`}
</CodeCard>
```

#### Callout

```mdx
<Callout title="Important" type="warn">
  This is a warning message.
</Callout>
```

Types: `info` (default), `warn`, `success`

#### Compare (Before/After)

```mdx
<Compare>
  <CompareCol kind="bad" pill="avoid">
    <p>Old approach with error handling</p>
    <code>value, err := DoSomething()</code>
  </CompareCol>
  <CompareCol kind="good" pill="recommended">
    <p>Functional approach with Result</p>
    <code>result := R.TryCatch(DoSomething)</code>
  </CompareCol>
</Compare>
```

#### Benchmark

```mdx
<Bench
  title="Performance comparison"
  command="go test -bench=. -benchmem"
  rows={[
    {label: 'Baseline', bar: 1, barKind: 'lose', nsOp: '14,820', bOp: '96'},
    {label: 'Optimized', bar: 0.14, barKind: 'win', nsOp: '2,140', bOp: '48', winner: true},
  ]}
/>
```

See [`docs/design-kit.mdx`](./docs/design-kit.mdx) for complete component documentation.

## 🎨 Customization

### Theme

Customize colors, fonts, and styling in [`src/css/custom.css`](./src/css/custom.css).

### Navigation

Update the navbar and footer in [`docusaurus.config.ts`](./docusaurus.config.ts).

### Sidebar

Modify sidebar structure in [`sidebars.ts`](./sidebars.ts).

## 🔍 Search

The site is configured for Algolia DocSearch. To enable search:

1. Apply for DocSearch at https://docsearch.algolia.com/apply/
2. Update credentials in `docusaurus.config.ts`:

```typescript
algolia: {
  appId: 'YOUR_APP_ID',
  apiKey: 'YOUR_API_KEY',
  indexName: 'fp-go',
}
```

See [`ALGOLIA_SETUP.md`](./ALGOLIA_SETUP.md) for detailed setup instructions.

## 📦 Available Scripts

| Script | Description |
|--------|-------------|
| `npm start` | Start development server |
| `npm run build` | Build for production |
| `npm run serve` | Serve production build locally |
| `npm run clear` | Clear Docusaurus cache |
| `npm run typecheck` | Run TypeScript type checking |
| `npm run swizzle` | Eject Docusaurus components for customization |
| `npm run deploy` | Deploy to GitHub Pages |

## 🚢 Deployment

### GitHub Pages

The site is configured for GitHub Pages deployment:

```bash
npm run build
npm run deploy
```

This will build and push to the `gh-pages` branch.

### Other Platforms

The `build/` directory can be deployed to any static hosting:

- **Netlify**: Connect your repo and set build command to `npm run build`
- **Vercel**: Import project and set build command to `npm run build`
- **AWS S3**: Upload `build/` contents to S3 bucket
- **Cloudflare Pages**: Connect repo with build command `npm run build`

## 📚 Documentation Sections

### Core Documentation
- **Getting Started**: Introduction, installation, quick start
- **Concepts**: Core FP concepts (monads, composition, pure functions)
- **Recipes**: Practical examples and patterns
- **Migration**: Guides for upgrading and interoperability

### API Reference
- **v2 (Current)**: Complete API documentation for v2.x
- **v1 (Legacy)**: Documentation for v1.x (deprecated)

### Advanced Topics
- **Architecture**: System design patterns
- **Performance**: Optimization techniques
- **Type Theory**: Advanced type system concepts

## 🛠️ Development Tips

### Hot Reload

Changes to `.md`, `.mdx`, and `.tsx` files trigger automatic reload.

### TypeScript

The project uses TypeScript for type safety. Run type checking:

```bash
npm run typecheck
```

### Component Development

Custom components are in `src/components/`. They're automatically available in MDX files via `src/theme/MDXComponents.tsx`.

### Debugging

Enable debug mode:

```bash
DEBUG=* npm start
```

### Clear Cache

If you encounter build issues:

```bash
npm run clear
npm start
```

## 🤝 Contributing

### Adding Documentation

1. Create your documentation file in the appropriate `docs/` subdirectory
2. Use frontmatter for metadata
3. Leverage custom components for rich content
4. Update `sidebars.ts` if needed
5. Test locally with `npm start`
6. Submit a pull request

### Style Guide

- Use clear, concise language
- Include code examples for concepts
- Add type signatures for functions
- Use custom components for callouts and comparisons
- Keep examples practical and runnable
- Add links to related documentation

### Code Examples

- Prefer complete, runnable examples
- Include imports when relevant
- Show both input and output
- Highlight key concepts
- Use `CodeCard` component with `status="tested"` for verified examples

## 📄 License

This documentation is part of the fp-go project and is licensed under the Apache License 2.0.

## 🔗 Links

- **Main Repository**: https://github.com/IBM/fp-go
- **Documentation Site**: https://ibm.github.io/fp-go/
- **Go Package**: https://pkg.go.dev/github.com/IBM/fp-go/v2
- **Playground Sandbox**: https://github.com/teyouale/fp-go-sandbox
- **Issues**: https://github.com/IBM/fp-go/issues
- **Discussions**: https://github.com/IBM/fp-go/discussions

## 📞 Support

- **Documentation Issues**: Open an issue in the main repository
- **Docusaurus Questions**: Check [Docusaurus documentation](https://docusaurus.io/)
- **Algolia Search**: See [ALGOLIA_SETUP.md](./ALGOLIA_SETUP.md)

