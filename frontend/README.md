# Leviosa Frontend

SvelteKit 5 application (TypeScript, Tailwind CSS v4, Node.js adapter). See `CLAUDE.md` for architecture details (route structure, auth flow, key libraries).

## Requirements

- Node.js 20+
- pnpm

## Development

```bash
pnpm install
pnpm run dev          # start dev server
pnpm run dev -- --open
```

Configure environment variables in `.env` (`API_URL`, `SESSION_COOKIE_NAME`, `CLIENT_IP_HEADER`, etc.) — see `CLAUDE.md` for the full list. In development and staging, auth falls back to a mock user (`$lib/data/user`); production validates real sessions against the backend.

## Building

```bash
pnpm run build
pnpm run preview   # preview the production build locally
```

## Type checking

```bash
pnpm run check
pnpm run check:watch
```

## Testing

No test framework is currently configured.
