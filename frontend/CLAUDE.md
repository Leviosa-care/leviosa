# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a SvelteKit 5 frontend application using TypeScript, Tailwind CSS v4, and the Node.js adapter. The project follows a modern SvelteKit structure with server-side authentication handling.

### Key Architecture Components

- **Authentication Flow**: Uses session-based authentication with cookies, handled in `hooks.server.ts`
- **Route Structure**: 
  - `src/routes/(app)/` - Protected application routes
  - `src/routes/auth/` - Authentication routes (sign in/up)
  - `src/routes/legal/` - Legal pages (public)
- **Environment-based Authentication**: Mock user in dev/staging, real authentication in production
- **API Integration**: Backend communication through fetch with automatic Bearer token injection

### File Structure

```
src/
├── hooks.server.ts          # Server-side request handling & auth
├── app.html                 # HTML template
├── app.css                  # Global styles
├── lib/                     # Reusable components and utilities
│   ├── ui/                  # UI components
│   ├── utils/               # Utility functions
│   ├── context/             # Svelte contexts
│   ├── data/                # Data models and mock data
│   ├── constructor/         # Constructor functions
│   ├── navigation/          # Navigation components
│   └── types/               # TypeScript type definitions
├── routes/                  # SvelteKit file-based routing
└── assets/                  # Static assets (SVG icons, etc.)
```

## Development Commands

```bash
# Development server
npm run dev
npm run dev -- --open  # Open in browser

# Build for production
npm run build

# Preview production build
npm run preview

# Type checking
npm run check
npm run check:watch  # Watch mode

# SvelteKit sync (prepare types)
npm run prepare
```

## Configuration Files

- `svelte.config.js` - SvelteKit configuration with Node.js adapter and path aliases
- `vite.config.ts` - Vite configuration with Tailwind CSS v4 and SvelteKitPWA (Workbox) integration
- `tailwind.config.ts` - Tailwind CSS v4 configuration
- `tsconfig.json` - TypeScript configuration extending SvelteKit defaults

## Path Aliases

- `@/*` → `$lib` (SvelteKit's lib alias)
- `@assets/*` → `src/assets`

## Authentication System

The app uses a session-based authentication system:

1. **Development/Staging**: Uses `mockUser` from `$lib/data/user`
2. **Production**: Validates sessions via API calls to `/users/me`
3. **Protected Routes**: All routes except `/auth` and `/legal` require authentication
4. **Session Handling**: Automatic token injection in fetch requests via `hooks.server.ts`

## Environment Variables

Located in `.env` file (private environment variables in SvelteKit):
- `NODE_ENV` - Environment mode
- `SESSION_COOKIE_NAME` - Session cookie identifier
- `API_URL` - Backend API URL
- `CLIENT_IP_HEADER` - Header for client IP forwarding

## Key Libraries

- **SvelteKit 5** - Full-stack web framework
- **Tailwind CSS v4** - Utility-first CSS framework
- **bits-ui** - Headless UI components for Svelte
- **sveltekit-superforms** - Form validation and handling
- **arktype** - TypeScript-first schema validation
- **@internationalized/date** - Date internationalization

## Testing

No test framework is currently configured. Check with the team for testing requirements before adding tests.
