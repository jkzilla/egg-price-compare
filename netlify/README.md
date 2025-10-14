# Netlify Deployment

Deploy Egg Price Comparison API to Netlify Functions.

## Prerequisites

1. Netlify account
2. Netlify CLI: `npm install -g netlify-cli`

## Setup

### 1. Install Netlify CLI

```bash
npm install -g netlify-cli
```

### 2. Login to Netlify

```bash
netlify login
```

### 3. Initialize Site

```bash
netlify init
```

### 4. Set Environment Variables

```bash
netlify env:set WALMART_API_KEY your_walmart_api_key
netlify env:set WALGREENS_API_KEY your_walgreens_api_key
```

## Local Development

```bash
# Install dependencies
netlify dev

# Access at http://localhost:8888
```

## Deploy

### Deploy to Production

```bash
netlify deploy --prod
```

### Deploy Preview

```bash
netlify deploy
```

## Configuration

The `netlify.toml` file configures:

- **Build command**: Compiles Go binary
- **Functions**: GraphQL endpoint at `/.netlify/functions/graphql`
- **Redirects**: Routes `/graphql` to the function
- **CORS**: Allows cross-origin requests

## Environment Variables

Set in Netlify dashboard or via CLI:

- `WALMART_API_KEY` - Walmart API key
- `WALGREENS_API_KEY` - Walgreens API key
- `GO_VERSION` - Go version (default: 1.21)

## Endpoints

After deployment:

- **GraphQL API**: `https://your-site.netlify.app/graphql`
- **GraphQL Playground**: `https://your-site.netlify.app/`

## Troubleshooting

### Build fails

Check build logs in Netlify dashboard:
```bash
netlify build
```

### Function errors

View function logs:
```bash
netlify functions:log graphql
```

### CORS issues

Verify headers in `netlify.toml`:
```toml
[[headers]]
  for = "/graphql"
  [headers.values]
    Access-Control-Allow-Origin = "*"
```

## Resources

- Netlify Docs: https://docs.netlify.com
- Netlify Functions: https://docs.netlify.com/functions/overview/
- Go on Netlify: https://docs.netlify.com/configure-builds/available-software/#go
