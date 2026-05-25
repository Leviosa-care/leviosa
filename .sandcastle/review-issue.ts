import { createSandbox, claudeCode } from "@ai-hero/sandcastle";
import { docker } from "@ai-hero/sandcastle/sandboxes/docker";
import { readFileSync, readdirSync } from "fs";
import { join } from "path";

const issueNum = process.argv[2]?.padStart(2, "0");

if (!issueNum) {
  console.error("Usage: npx tsx .sandcastle/review-issue.ts <issue-number>");
  console.error("Example: npx tsx .sandcastle/review-issue.ts 1");
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

const reviewAgent = claudeCode("claude-sonnet-4-6");

console.log(`\nReviewing issue ${issueNum} on branch ${branchName}...\n`);

await using sandbox = await createSandbox({
  branch: branchName,
  baseBranch: "main",
  sandbox: docker({
    mounts: [
      { hostPath: "~/.claude", sandboxPath: "/home/agent/.claude" },
      { hostPath: "~/.npm", sandboxPath: "/home/agent/.npm", readonly: true },
    ],
  }),
  hooks: {
    sandbox: {
      onSandboxReady: [
        { command: "cd frontend && npm install --prefer-offline", timeoutMs: 180_000 },
        { command: "cd backend && go mod download", timeoutMs: 120_000 },
      ],
    },
  },
});

await sandbox.run({
  agent: reviewAgent,
  name: `issue-${issueNum}-review`,
  maxIterations: 10,
  prompt: `Review the current state of the branch against the issue below and fix any problems.

## Issue

${issueContent}

## Review checklist

1. Verify every acceptance criterion is satisfied — read the code, do not assume
2. Run \`cd backend && go build ./...\` — fix any compilation errors
3. Run \`cd frontend && npm run check\` — fix any TypeScript errors
4. Check that no getMock* functions or hardcoded mock data remain in +page.server.ts files
5. Check that session cookies are forwarded correctly on all server-side API calls
6. Verify error handling exists at API call boundaries (401 redirect, 500 message)
7. Remove any TODO markers or leftover scaffolding comments from the issue text

Make only targeted fixes. Do not refactor working code or expand scope. Stop when all criteria are met and verification passes.`,
});

console.log(`\nReview complete. Branch: ${branchName}`);
console.log(`Inspect with: git diff main...${branchName}\n`);
