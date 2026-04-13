---
title: Common Issues
description: Frequently encountered issues when working with GitHub Agentic Workflows and their solutions.
sidebar:
  order: 200
---

This reference documents frequently encountered issues when working with GitHub Agentic Workflows, organized by workflow stage and component.

## Installation Issues

### Extension Installation Fails

If `gh extension install github/gh-aw` fails, use the standalone installer (works in Codespaces and restricted networks):

```bash wrap
curl -sL https://raw.githubusercontent.com/github/gh-aw/main/install-gh-aw.sh | bash
```

For specific versions, pass the tag as an argument ([see releases](https://github.com/github/gh-aw/releases)):

```bash wrap
curl -sL https://raw.githubusercontent.com/github/gh-aw/main/install-gh-aw.sh | bash -s -- v0.40.0
```

Verify with `gh extension list`.

## Organization Policy Issues

### Custom Actions Not Allowed in Enterprise Organizations

**Error Message:**

```text
The action github/gh-aw/actions/setup@a933c835b5e2d12ae4dead665a0fdba420a2d421 is not allowed in {ORG} because all actions must be from a repository owned by your enterprise, created by GitHub, or verified in the GitHub Marketplace.
```

**Cause:** Enterprise policies restrict which GitHub Actions can be used. Workflows use `github/gh-aw/actions/setup` which may not be allowed.

**Solution:** Enterprise administrators must allow `github/gh-aw` in the organization's action policies:

#### Option 1: Allow Specific Repositories (Recommended)

Add `github/gh-aw` to your organization's allowed actions list:

1. Navigate to your organization's settings: `https://github.com/organizations/YOUR_ORG/settings/actions`
2. Under **Policies**, select **Allow select actions and reusable workflows**
3. In the **Allow specified actions and reusable workflows** section, add:
   ```text
   github/gh-aw@*
   ```
4. Save the changes

See GitHub's docs on [managing Actions permissions](https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization#allowing-select-actions-and-reusable-workflows-to-run).

#### Option 2: Configure Organization-Wide Policy File

Add `github/gh-aw@*` to your centralized `policies/actions.yml` and commit to your organization's `.github` repository. See GitHub's docs on [community health files](https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/creating-a-default-community-health-file).

```yaml
allowed_actions:
  - "actions/*"
  - "github/gh-aw@*"
```

#### Verification

Wait a few minutes for policy propagation, then re-run your workflow. If issues persist, verify at `https://github.com/organizations/YOUR_ORG/settings/actions`.

> [!TIP]
> The gh-aw actions are open source at [github.com/github/gh-aw/tree/main/actions](https://github.com/github/gh-aw/tree/main/actions) and pinned to specific SHAs for security.

## Repository Configuration Issues

### Actions Restrictions Reported During Init

The CLI validates three permission layers. Fix restrictions in Repository Settings → Actions → General:

1. **Actions disabled**: Enable Actions ([docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository))
2. **Local-only**: Switch to "Allow all actions" or enable GitHub-created actions ([docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#managing-github-actions-permissions-for-your-repository))
3. **Selective allowlist**: Enable "Allow actions created by GitHub" checkbox ([docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#allowing-select-actions-and-reusable-workflows-to-run))

> [!NOTE]
> Organization policies override repository settings. Contact admins if settings are grayed out.

## Workflow Compilation Issues

### Frontmatter Field Not Taking Effect

If a frontmatter setting appears to be silently ignored, the field name may be misspelled. The compiler does not warn about unknown field names — they are silently discarded.

> [!WARNING]
> Common frontmatter field name mistakes:
>
> | Wrong | Correct |
> |-------|---------|
> | `agent:` | `engine:` |
> | `mcp-servers:` | `tools:` (under which MCP servers are configured) |
> | `tool-sets:` | `toolsets:` (under `tools.github:`) |
> | `allowed_repos:` | `allowed-repos:` (under `tools.github:`) |
> | `timeout:` | `timeout-minutes:` |
>
> Run `gh aw compile --verbose` to confirm which settings were parsed. If your setting is missing from the output, check the [Frontmatter Reference](/gh-aw/reference/frontmatter/) for the correct field name.

### Workflow Won't Compile

Check YAML frontmatter syntax (indentation, colons with spaces), verify required fields (`on:`), and ensure types match the schema. Use `gh aw compile --verbose` for details.

### Lock File Not Generated

Fix compilation errors (`gh aw compile 2>&1 | grep -i error`) and verify write permissions on `.github/workflows/`.

### Orphaned Lock Files

Remove old `.lock.yml` files with `gh aw compile --purge` after deleting `.md` workflow files.

## Import and Include Issues

### Import File Not Found

Import paths are relative to repository root. Verify with `git status` (e.g., `.github/workflows/shared/tools.md`).

### Multiple Agent Files Error

Import only one `.github/agents/` file per workflow.

### Circular Import Dependencies

Compilation hangs indicate circular imports. Remove circular references.

## Tool Configuration Issues

### GitHub Tools Not Available

Configure using `toolsets:` ([tools reference](/gh-aw/reference/github-tools/)):

```yaml wrap
tools:
  github:
    toolsets: [repos, issues]
```

### Toolset Missing Expected Tools

Check [GitHub Toolsets](/gh-aw/reference/github-tools/), combine toolsets (`toolsets: [default, actions]`), or inspect with `gh aw mcp inspect <workflow>`.

### MCP Server Connection Failures

Verify package installation, syntax, and environment variables:

```yaml
mcp-servers:
  my-server:
    command: "npx"
    args: ["@myorg/mcp-server"]
    env:
      API_KEY: "${{ secrets.MCP_API_KEY }}"
```

### Playwright Network Access Denied

Add domains to `network.allowed`:

```yaml wrap
network:
  allowed:
    - github.com
    - "*.github.io"
```

### Cannot Find Module 'playwright'

**Error:**

```text
Error: Cannot find module 'playwright'
```

**Cause:** The agent tried to `require('playwright')` but Playwright is provided through MCP tools, not as an npm package.

**Solution:** Use MCP Playwright tools:

```javascript
// ❌ INCORRECT - This won't work
const playwright = require('playwright');
const browser = await playwright.chromium.launch();

// ✅ CORRECT - Use MCP Playwright tools
// Example: Navigate and take screenshot
await mcp__playwright__browser_navigate({
  url: "https://example.com"
});

await mcp__playwright__browser_snapshot();

// Example: Execute custom Playwright code
await mcp__playwright__browser_run_code({
  code: `async (page) => {
    await page.setViewportSize({ width: 390, height: 844 });
    const title = await page.title();
    return { title, url: page.url() };
  }`
});
```

See [Playwright Tool documentation](/gh-aw/reference/tools/#playwright-tool-playwright) for all available tools.

### Playwright MCP Initialization Failure (EOF Error)

**Error:**

```text
Failed to register tools error="initialize: EOF" name=playwright
```

**Cause:** Chromium crashes before tool registration completes due to missing Docker security flags.

**Solution:** Upgrade to version 0.41.0+ which includes required Docker flags:

```bash wrap
gh extension upgrade gh-aw
```

## Permission Issues

### Write Operations Fail

Agentic workflows cannot write to GitHub directly. All writes (issues, comments, PR updates)
must go through the `safe-outputs` system, which validates and executes write operations on
behalf of the workflow.

Ensure your workflow frontmatter declares the safe output types it needs:

```yaml wrap
safe-outputs:
  create-issue:
    title-prefix: "[bot] "
    labels: [automation]
  add-comment:      # no configuration required; uses defaults
  update-issue:     # no configuration required; uses defaults
```

If the operation you need is not listed in the [Safe Outputs reference](/gh-aw/reference/safe-outputs/),
it may not be supported yet. See the [Safe Outputs Specification](/gh-aw/reference/safe-outputs-specification/)
for the full list of available output types and their configuration options.

### Safe Outputs Not Creating Issues

Disable staged mode:

```yaml wrap
safe-outputs:
  staged: false
  create-issue:
    title-prefix: "[bot] "
    labels: [automation]
```

### Project Field Type Errors

GitHub Projects reserves field names like `REPOSITORY`. Use alternatives (`repo`, `source_repository`, `linked_repo`):

```yaml wrap
# ❌ Wrong: repository
# ✅ Correct: repo
safe-outputs:
  update-project:
    fields:
      repo: "myorg/myrepo"
```

Delete conflicting fields in Projects UI and recreate.

## Engine-Specific Issues

### Copilot CLI Not Found

Verify compilation succeeded. Compiled workflows include CLI installation steps.

### Model Not Available

Use default (`engine: copilot`) or specify available model (`engine: {id: copilot, model: gpt-4}`).

### Copilot License or Inference Access Issues

If the Copilot inference step fails despite `COPILOT_GITHUB_TOKEN` being configured, the PAT owner may lack a valid Copilot license. Test locally:

```bash
export COPILOT_GITHUB_TOKEN="<your-github-pat>"
copilot -p "write a haiku"
```

If this fails, contact your organization administrator. The token owner must have an active Copilot subscription — organization-managed licenses may restrict programmatic API access.

## GitHub Enterprise Server Issues

> [!TIP]
> For a complete walkthrough of setting up and debugging workflows on **GHE Cloud with data residency** (`*.ghe.com`), see [Debugging GHE Cloud with Data Residency](/gh-aw/troubleshooting/debug-ghe/).

### Copilot Engine Prerequisites on GHES

Before running Copilot-based workflows on GHES, verify the following:

**Site admin:** Enable GitHub Connect; purchase and activate Copilot licensing at the enterprise level; allow outbound HTTPS to `api.githubcopilot.com` and `api.enterprise.githubcopilot.com`.

**Enterprise/org admin:** Assign Copilot seats to the `COPILOT_GITHUB_TOKEN` owner; enable Copilot usage in the organization's Copilot policy.

**Workflow configuration:**

```aw wrap
engine:
  id: copilot
  api-target: api.enterprise.githubcopilot.com
network:
  allowed:
    - defaults
    - api.enterprise.githubcopilot.com
```

See [Enterprise API Endpoint](/gh-aw/reference/engines/#enterprise-api-endpoint-api-target) for GHEC/GHES `api-target` values.

### Copilot GHES: Common Error Messages

**`Error loading models: 400 Bad Request`**

Copilot is not licensed at the enterprise level or the API proxy is routing incorrectly. Verify enterprise Copilot settings and confirm GitHub Connect is enabled.

**`403 "unauthorized: not licensed to use Copilot"`**

No Copilot license or seat is assigned to the PAT owner. Contact the site admin to enable Copilot at the enterprise level, then have an org admin assign a seat to the token owner.

**`403 "Resource not accessible by personal access token"`**

Wrong token type or missing permissions. Use a fine-grained PAT with the **Copilot Requests: Read** account permission, or a classic PAT with the `copilot` scope. See [`COPILOT_GITHUB_TOKEN`](/gh-aw/reference/auth/#copilot_github_token) for setup instructions.

**`Could not resolve to a Repository`**

`GH_HOST` is not set when running `gh` commands. This typically occurs in custom frontmatter jobs from older compiled workflows. Recompile with `gh aw compile` — compiled workflows now automatically export `GH_HOST` in custom jobs.

For local CLI commands (`gh aw audit`, `gh aw add-wizard`), ensure you are inside a GHES repository clone or set `GH_HOST` explicitly:

```bash wrap
GH_HOST=github.company.com gh aw audit <run-id>
```

**Firewall blocks outbound HTTPS to `api.<ghes-host>`**

Add the GHES domain to your workflow's allowed list:

```aw wrap
engine:
  id: copilot
  api-target: api.company.ghe.com
network:
  allowed:
    - defaults
    - company.ghe.com
    - api.company.ghe.com
```

**`gh aw add-wizard` or `gh aw init` creates a PR on github.com instead of GHES**

Run these commands from inside a GHES repository clone — they auto-detect the GHES host from the git remote. If PR creation still fails, use `gh aw add` to generate the workflow file, then create the PR manually with `gh pr create`.

## Context Expression Issues

### Unauthorized Expression

Use only [allowed expressions](/gh-aw/reference/templating/) (`github.event.issue.number`, `github.repository`, `steps.sanitized.outputs.text`). Disallowed: `secrets.*`, `env.*`.

### Sanitized Context Empty

`steps.sanitized.outputs.text` requires issue/PR/comment events (`on: issues:`), not `push:` or similar triggers.

## Build and Test Issues

### Documentation Build Fails

Clean install and rebuild:

```bash wrap
cd docs
rm -rf node_modules package-lock.json
npm install
npm run build
```

Check for malformed frontmatter, MDX syntax errors, or broken links.

### Tests Failing After Changes

Format and lint before testing:

```bash wrap
make fmt
make lint
make test-unit
```

## Network and Connectivity Issues

### Firewall Denials for Package Registries

Add ecosystem identifiers ([Network Configuration Guide](/gh-aw/guides/network-configuration/)):

```yaml wrap
network:
  allowed:
    - defaults    # Infrastructure
    - python      # PyPI
    - node        # npm
    - containers  # Docker
    - go          # Go modules
```

### URLs Appearing as "(redacted)"

Add domains to allowed list ([Network Permissions](/gh-aw/reference/network/)):

```yaml wrap
network:
  allowed:
    - defaults
    - "api.example.com"
```

### Cannot Download Remote Imports

Verify network (`curl -I https://raw.githubusercontent.com/github/gh-aw/main/README.md`) and auth (`gh auth status`).

### MCP Server Connection Timeout

Use local servers (`command: "node"`, `args: ["./server.js"]`).

## Cache Issues

### Cache Not Restoring

Verify key patterns match (caches expire after 7 days):

```yaml wrap
cache:
  key: deps-${{ hashFiles('package-lock.json') }}
  restore-keys: deps-
```

### Cache Memory Not Persisting

Configure cache for memory MCP server:

```yaml wrap
tools:
  cache-memory:
    key: memory-${{ github.workflow }}-${{ github.run_id }}
```

## Integrity Filtering Blocking Expected Content

In public repositories, `min-integrity: approved` is applied automatically, restricting agent visibility to owners, members, and collaborators. Workflows may miss issues, PRs, or comments from external contributors.

**Option 1: Keep the default level (Recommended)** — For sensitive operations (code generation, repository updates, web access), use separate workflows, manual triggers, or approval stages.

**Option 2: Lower the integrity level** — Only if your workflow validates input, uses restrictive safe outputs, and doesn't access secrets:

```yaml wrap
tools:
  github:
    min-integrity: none
```

For community triage that needs contributor input but not anonymous users, `min-integrity: unapproved` is a useful middle ground.

See [Integrity Filtering](/gh-aw/reference/integrity/) for details.

## Workflow Failures and Debugging

### Workflow Job Timed Out

The default timeout is 20 minutes. Increase it in workflow frontmatter, then recompile with `gh aw compile`. If timeouts persist, reduce the task scope or split into smaller workflows. See [Long Build Times](/gh-aw/reference/sandbox/#long-build-times) for caching strategies and self-hosted runner recommendations.

### Engine Timeout Error Messages

| Engine | Error pattern | Setting |
|--------|--------------|---------|
| GitHub Actions | `exceeded the maximum execution time of N minutes` | `timeout-minutes` |
| Claude (tool call) | `Bash tool timed out after 60 seconds` | `tools.timeout` |
| Claude (max turns) | `Reached maximum number of turns (N). Stopping.` | `max-turns` |
| Codex | `Tool call timed out after 120 seconds` | `tools.timeout` |
| MCP server startup | `Failed to register tools error="initialize: timeout"` | `tools.startup-timeout` |
| Copilot autopilot | *(task silently incomplete when `max-continuations` exhausted)* | `max-continuations` |

Common timeout settings:

```yaml wrap
timeout-minutes: 60      # GitHub Actions job limit (default: 20)
max-turns: 30            # Claude turns limit
max-continuations: 5     # Copilot autopilot run limit
tools:
  timeout: 600           # Per-tool-call limit in seconds (default: 60–120)
  startup-timeout: 300   # MCP server startup budget (default: 120)
```

### Debugging a Failing Workflow

Common causes: missing tokens, permission mismatches, network restrictions, disabled tools, or rate limits. The fastest way to debug is to ask an agent — load the `agentic-workflows` agent and give it the run URL:

**Using Copilot Chat** (requires [agentic authoring setup](/gh-aw/guides/agentic-authoring/#configuring-your-repository)):

```text wrap
/agent agentic-workflows debug https://github.com/OWNER/REPO/actions/runs/RUN_ID
```

**Using any coding agent** (self-contained, no setup required):

```text wrap
Debug this workflow run using https://raw.githubusercontent.com/github/gh-aw/main/debug.md

The failed workflow run is at https://github.com/OWNER/REPO/actions/runs/RUN_ID
```

You can also investigate manually: check logs (`gh aw logs`), audit the run (`gh aw audit <run-id>`), enable verbose mode (`--verbose`), set `ACTIONS_STEP_DEBUG = true`, or inspect MCP config (`gh aw mcp inspect`). See the [Debugging Workflows](/gh-aw/troubleshooting/debugging/) guide for a comprehensive walkthrough.

### Enable Debug Logging

The `DEBUG` environment variable activates detailed internal logging for any `gh aw` command:

```bash
DEBUG=* gh aw compile                              # all logs
DEBUG=workflow:* gh aw compile my-workflow         # specific package
DEBUG=workflow:*,cli:* gh aw compile my-workflow   # multiple packages
DEBUG=*,-workflow:test gh aw compile my-workflow   # exclude a logger
DEBUG_COLORS=0 DEBUG=* gh aw compile 2>&1 | tee debug.log  # capture to file
```

Debug output goes to `stderr`. Each log line shows the namespace (`workflow:compiler`), message, and time elapsed since the previous entry. Common namespaces: `cli:compile_command`, `workflow:compiler`, `workflow:expression_extraction`, `parser:frontmatter`. Wildcards match any suffix (`workflow:*`).

## Operational Runbooks

See [Workflow Health Monitoring Runbook](https://github.com/github/gh-aw/blob/main/.github/aw/runbooks/workflow-health.md) for diagnosing errors.

## Getting Help

Review [reference docs](/gh-aw/reference/workflow-structure/), search [existing issues](https://github.com/github/gh-aw/issues), or create an issue. See [Error Reference](/gh-aw/troubleshooting/errors/) and [Frontmatter Reference](/gh-aw/reference/frontmatter/).
