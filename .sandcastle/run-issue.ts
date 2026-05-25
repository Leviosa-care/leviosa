import { createSandbox, claudeCode } from "@ai-hero/sandcastle";
import { docker } from "@ai-hero/sandcastle/sandboxes/docker";
import { readFileSync, readdirSync } from "fs";
import { join } from "path";

// ---------------------------------------------------------------------------
// Validate env
// ---------------------------------------------------------------------------

const requiredEnv = ["ZAI_API_KEY", "ZAI_BASE_URL", "GLM_SONNET_MODEL", "GLM_OPUS_MODEL", "GLM_HAIKU_MODEL"];
const missing = requiredEnv.filter((k) => !process.env[k]);
if (missing.length > 0) {
  console.error(`Missing env vars: ${missing.join(", ")}`);
  console.error("These are defined in ~/.profile — make sure you're in a login shell (source ~/.profile).");
  process.exit(1);
}

// ---------------------------------------------------------------------------
// Resolve issue file
// ---------------------------------------------------------------------------

const issueNum = process.argv[2]?.padStart(2, "0");

if (!issueNum) {
  console.error("Usage: npx tsx .sandcastle/run-issue.ts <issue-number>");
  console.error("Example: npx tsx .sandcastle/run-issue.ts 1");
  process.exit(1);
}

const issuesDir = join(process.cwd(), "docs/issues");
const issueFile = readdirSync(issuesDir).find((f) => f.startsWith(`${issueNum}-`));

if (!issueFile) {
  console.error(`No issue file found starting with '${issueNum}-' in docs/issues/`);
  process.exit(1);
}

const issueContent = readFileSync(join(issuesDir, issueFile), "utf-8");
const branchName = `agent/issue-${issueNum}`;

// ---------------------------------------------------------------------------
// Agents
// ---------------------------------------------------------------------------

// GLM via Z.ai — mirrors the `cg` shell alias env var pattern
const glmAgent = claudeCode("claude-sonnet-4-6", {
  effort: "high",
  env: {
    ANTHROPIC_AUTH_TOKEN: process.env.ZAI_API_KEY!,
    ANTHROPIC_BASE_URL: process.env.ZAI_BASE_URL!,
    ANTHROPIC_DEFAULT_OPUS_MODEL: process.env.GLM_OPUS_MODEL!,
    ANTHROPIC_DEFAULT_SONNET_MODEL: process.env.GLM_SONNET_MODEL!,
    ANTHROPIC_DEFAULT_HAIKU_MODEL: process.env.GLM_HAIKU_MODEL!,
    API_TIMEOUT_MS: "3000000",
  },
});

// Standard Claude Code (Anthropic) — review pass
const reviewAgent = claudeCode("claude-sonnet-4-6");

// ---------------------------------------------------------------------------
// Pipeline
// ---------------------------------------------------------------------------

console.log(`\nStarting pipeline for issue ${issueNum} → branch ${branchName}\n`);

await using sandbox = await createSandbox({
  branch: branchName,
  baseBranch: "main",
  sandbox: docker({
    mounts: [
      // Claude Code credentials — required for both agents (GLM overrides ANTHROPIC_AUTH_TOKEN on top).
      // Must NOT be readonly: Claude Code writes session state here; EROFS breaks Bash and commits.
      { hostPath: "~/.claude", sandboxPath: "/home/agent/.claude" },
      // npm cache — avoids re-downloading on every run
      { hostPath: "~/.npm", sandboxPath: "/home/agent/.npm", readonly: true },
    ],
  }),
  // Pre-warm dependencies so agents can run build/check commands immediately
  hooks: {
    sandbox: {
      onSandboxReady: [
        { command: "cd frontend && npm install --prefer-offline", timeoutMs: 180_000 },
        { command: "cd backend && go mod download", timeoutMs: 120_000 },
      ],
    },
  },
});

// ---------------------------------------------------------------------------
// Step 1: GLM implements
// ---------------------------------------------------------------------------

console.log(`[1/2] GLM (${process.env.GLM_SONNET_MODEL}) implementing...\n`);

await sandbox.run({
  agent: glmAgent,
  name: `issue-${issueNum}-implement`,
  maxIterations: 20,
  completionSignal: "<promise>COMPLETE</promise>",
  prompt: `You are implementing a specific issue for the Leviosa wellness booking platform.

This is a Go 1.24 + SvelteKit 5 monorepo. The backend follows hexagonal architecture. The frontend uses TypeScript strict mode.

## Issue

${issueContent}

## Conventions

**Backend (Go):**
- Hexagonal layers: domain → ports → application → interface (HTTP handlers) / infrastructure (adapters)
- New service methods go on the port interface first, then the application layer implements them
- New HTTP handlers follow the existing pattern in the same service's interface directory
- Fields with PII are tagged \`encx:"encrypt"\` — the repository layer handles this transparently
- No code comments unless the WHY is non-obvious

**Frontend (SvelteKit 5):**
- Use \`locals.user.id\` for the authenticated user's ID in +page.server.ts
- Forward the session cookie on all server-side fetches:
  \`headers: { Cookie: \`\${locals.sessionCookieName}=\${cookies.get(locals.sessionCookieName)}\` }\`
- Forms: sveltekit-superforms + arktype, following existing page patterns
- Remove getMock* functions entirely — no mock data in the final result

**Verification (run these before finishing):**
- \`cd backend && go build ./...\` must succeed
- \`cd frontend && npm run check\` must succeed

When ALL acceptance criteria are met and both verification commands pass, output exactly:
<promise>COMPLETE</promise>`,
});

// ---------------------------------------------------------------------------
// Step 2: Claude Code reviews
// ---------------------------------------------------------------------------

console.log(`\n[2/2] Claude Code (Anthropic) reviewing...\n`);

try {
  await sandbox.run({
    agent: reviewAgent,
    name: `issue-${issueNum}-review`,
    maxIterations: 10,
    prompt: `The GLM model has just implemented an issue. Your job is to review the result and fix any problems you find.

## Issue that was implemented

${issueContent}

## Review checklist

1. Verify every acceptance criterion is satisfied — read the code, do not assume
2. Run \`cd backend && go build ./...\` — fix any compilation errors
3. Run \`cd frontend && npm run check\` — fix any TypeScript errors
4. Check that no getMock* functions or hardcoded mock data remain in +page.server.ts files
5. Check that session cookies are forwarded correctly on all server-side API calls
6. Verify error handling exists at API call boundaries (401 redirect, 500 message)
7. Remove any TODO markers or leftover scaffolding comments from the issue text

Make only targeted fixes. Do not refactor working code or expand the scope. Stop when all criteria are met and verification passes.`,
  });
  console.log(`\nPipeline complete. Branch: ${branchName}`);
} catch (err: unknown) {
  const msg = err instanceof Error ? err.message : String(err);
  if (msg.includes("limit") || msg.includes("rate")) {
    console.warn(`\n⚠️  Review step hit a usage limit — GLM implementation is on branch ${branchName}.`);
    console.warn(`Re-run the review later with: npm run review ${issueNum}\n`);
  } else {
    throw err;
  }
}

console.log(`Review with: git diff main...${branchName}\n`);
