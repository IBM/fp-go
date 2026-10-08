# Algolia DocSearch Setup Guide

This guide explains how to set up Algolia DocSearch for the fp-go documentation site.

## What is Algolia DocSearch?

Algolia DocSearch is a free search service for open-source documentation. It provides:
- Fast, typo-tolerant search
- Instant results as you type
- Keyboard navigation
- Mobile-friendly search UI
- Automatic indexing of your documentation

## Setup Steps

### 1. Apply for DocSearch

1. Go to https://docsearch.algolia.com/apply/
2. Fill out the application form with:
   - **Website URL**: Your deployed documentation URL (e.g., `https://ibm.github.io/fp-go/`)
   - **Email**: Your email address
   - **Repository**: `https://github.com/IBM/fp-go`
   - Confirm that:
     - You are the owner of the website
     - The website is publicly available
     - The website is an open-source project or technical documentation

3. Wait for approval (usually takes a few days)

### 2. Receive Credentials

Once approved, Algolia will send you:
- `appId`: Your application ID
- `apiKey`: Your public search API key
- `indexName`: Your index name (usually your project name)

### 3. Update Configuration

Update `docusaurus.config.ts` with your credentials:

```typescript
algolia: {
  appId: 'YOUR_APP_ID',           // Replace with your App ID
  apiKey: 'YOUR_SEARCH_API_KEY',  // Replace with your API Key
  indexName: 'fp-go',             // Replace with your index name
  contextualSearch: true,
  searchPagePath: 'search',
},
```

### 4. Deploy Your Site

Deploy your documentation site so Algolia can crawl it:

```bash
npm run build
# Deploy to your hosting (GitHub Pages, Netlify, Vercel, etc.)
```

### 5. Configure Crawler (Optional)

If you need custom crawler configuration, create a `.algolia/config.json` file:

```json
{
  "index_name": "fp-go",
  "start_urls": [
    "https://your-site.com/docs/"
  ],
  "sitemap_urls": [
    "https://your-site.com/sitemap.xml"
  ],
  "selectors": {
    "lvl0": {
      "selector": ".menu__link--sublist.menu__link--active",
      "global": true,
      "default_value": "Documentation"
    },
    "lvl1": "article h1",
    "lvl2": "article h2",
    "lvl3": "article h3",
    "lvl4": "article h4",
    "lvl5": "article h5",
    "text": "article p, article li"
  }
}
```

## Current Configuration

The fp-go documentation is already configured with Algolia DocSearch placeholders in `docusaurus.config.ts`.

### Features Enabled

- **Contextual Search**: Search results are filtered by the current documentation version
- **Search Page**: Dedicated search page at `/search`
- **Keyboard Shortcuts**: Press `/` or `Ctrl+K` to open search
- **Version-aware**: Searches respect the current documentation version (v1 vs v2)

## Testing Search Locally

Before Algolia is set up, Docusaurus provides a basic client-side search:

```bash
npm start
```

Then press `/` or click the search icon to test the search UI.

## Alternative: Local Search

If you prefer not to use Algolia, you can use local search plugins:

### Option 1: @docusaurus/theme-search-algolia (default)
Already configured, just needs Algolia credentials.

### Option 2: docusaurus-search-local
For fully local search without external services:

```bash
npm install @easyops-cn/docusaurus-search-local
```

Then update `docusaurus.config.ts`:

```typescript
themes: [
  [
    require.resolve("@easyops-cn/docusaurus-search-local"),
    {
      hashed: true,
      language: ["en"],
      highlightSearchTermsOnTargetPage: true,
      explicitSearchResultPath: true,
    },
  ],
],
```

## Search Best Practices

### 1. Use Descriptive Titles
```markdown
---
title: Error Handling Patterns
description: Common patterns for handling errors functionally
---
```

### 2. Add Keywords
Include relevant keywords in your content for better search results.

### 3. Structure Content
Use proper heading hierarchy (h1 → h2 → h3) for better indexing.

### 4. Add Metadata
Use frontmatter to add searchable metadata:

```markdown
---
title: Option Type
description: Handle optional values safely
keywords: [option, maybe, optional, null, undefined]
---
```

## Troubleshooting

### Search Not Working
1. Check that Algolia credentials are correct
2. Verify your site is deployed and accessible
3. Wait for Algolia to crawl your site (can take 24-48 hours after first deployment)
4. Check browser console for errors

### Search Results Outdated
1. Algolia crawls your site periodically (usually weekly)
2. You can request a manual re-crawl from the Algolia dashboard
3. Or set up a webhook to trigger crawls on deployment

### No Results Found
1. Ensure your content has proper HTML structure
2. Check that selectors in crawler config match your HTML
3. Verify content is not hidden behind authentication

## Resources

- [Algolia DocSearch Documentation](https://docsearch.algolia.com/docs/what-is-docsearch)
- [Docusaurus Search Documentation](https://docusaurus.io/docs/search)
- [Algolia Dashboard](https://www.algolia.com/dashboard)

## Support

For issues with:
- **Algolia setup**: Contact Algolia support or check their documentation
- **Docusaurus integration**: Check Docusaurus documentation or GitHub issues
- **fp-go documentation**: Open an issue at https://github.com/IBM/fp-go/issues