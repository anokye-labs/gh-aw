---
title: DispatchOps
description: Manually trigger and test agentic workflows with custom inputs using workflow_dispatch
sidebar:
  badge: { text: 'Manual', variant: 'tip' }
---

DispatchOps enables manual workflow execution via the GitHub Actions UI or CLI. Use it for research tasks, operational commands, testing, debugging, or any task that doesn't fit a schedule or event trigger.

## Trigger Configuration

Add `workflow_dispatch:` to your frontmatter, optionally with inputs:

```yaml
on:
  workflow_dispatch:
```

### With Input Parameters

```yaml
on:
  workflow_dispatch:
    inputs:
      topic:
        description: 'Research topic'
        required: true
        type: string
      priority:
        description: 'Task priority'
        required: false
        type: choice
        options:
          - low
          - medium
          - high
        default: medium
      deploy_target:
        description: 'Deployment environment'
        required: false
        type: environment
        default: staging
```

Supported input types: `string` (text), `boolean` (checkbox), `choice` (dropdown), `environment` (GitHub environments dropdown). The `environment` type auto-populates from repository Settings → Environments; it does not enforce protection rules — use `manual-approval:` for approval gates.

## Security

Triggering requires write access or higher. Restrict access with `roles:` or `bots:` fields:

```yaml
on:
  workflow_dispatch:
roles: [admin, maintainer]
bots: ["dependabot[bot]"]
```

Forks cannot trigger workflows in the parent repository — `workflow_dispatch` only executes in the repository where it's defined.

### Approval Gates

Require manual approval before execution using GitHub environment protection rules:

```yaml
on:
  workflow_dispatch:
manual-approval: production
```

Configure reviewers and wait timers in repository Settings → Environments. See [GitHub's environment documentation](https://docs.github.com/en/actions/deployment/targeting-different-environments/using-environments-for-deployment).

## Running Workflows

**From GitHub.com:** Go to the **Actions** tab, select the workflow, click **Run workflow**, fill in inputs, and confirm. Only workflows with `workflow_dispatch:` in their `on:` section appear — if yours is missing, verify it has been compiled and the `.lock.yml` pushed.

**From the CLI:** Use `gh aw run`, which matches by filename prefix and returns the run URL immediately:

```bash
gh aw run research --raw-field topic="quantum computing"

gh aw run scout \
  --raw-field topic="AI safety research" \
  --raw-field priority=high \
  --wait                         # Monitor and exit with success/failure code

gh aw run research --ref feature-branch    # Run from a specific branch
gh aw run workflow --repo owner/repository # Run in another repository
```

## Declaring and Referencing Inputs

### Referencing Inputs in Markdown

Access input values using GitHub Actions expression syntax:

```aw wrap
---
on:
  workflow_dispatch:
    inputs:
      topic:
        description: 'Research topic'
        required: true
        type: string
      depth:
        description: 'Analysis depth'
        type: choice
        options:
          - brief
          - detailed
        default: brief
permissions:
  contents: read
safe-outputs:
  create-discussion:
---

# Research Assistant

Research the following topic: "${{ github.event.inputs.topic }}"

Analysis depth requested: ${{ github.event.inputs.depth }}

Provide a ${{ github.event.inputs.depth }} analysis with key findings and recommendations.
```

Reference inputs with `${{ github.event.inputs.INPUT_NAME }}`. Use Handlebars conditionals to change behavior based on input values:

```markdown
{{#if (eq github.event.inputs.include_code "true")}}
Include actual code snippets in your analysis.
{{else}}
Describe code patterns without including actual code.
{{/if}}

{{#if (eq github.event.inputs.priority "high")}}
URGENT: Prioritize speed over completeness.
{{/if}}
```

## Branch Testing

Add `workflow_dispatch:` to feature branches for testing before merging:

```bash
gh aw trial ./research.md --raw-field topic="test query"  # isolated, no side effects
gh aw run research --ref feature/improve-workflow          # runs against live repo
```

## Common Use Cases

**On-demand research:** Add a `topic` string input and trigger with `gh aw run research --raw-field topic="AI safety"` when needed.

**Manual operations:** Use a `choice` input with predefined operations (cleanup, sync, audit) to execute specific tasks on demand.

**Testing and debugging:** Add `workflow_dispatch` to event-triggered workflows (issues, PRs) with optional test URL inputs to test without creating real events.

**Scheduled workflow testing:** Combine `schedule` with `workflow_dispatch` to test scheduled workflows immediately rather than waiting for the cron schedule.

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| Workflow not listed in GitHub UI | Verify `workflow_dispatch:` is in `on:`, compile (`gh aw compile workflow`), and push both `.md` and `.lock.yml` |
| "Workflow not found" | Use the filename without `.md` extension (e.g. `research` not `research.md`) |
| "Workflow cannot be run" | Add `workflow_dispatch:` to `on:`, recompile, verify `.lock.yml` includes the trigger |
| Permission denied | Verify write access; check the `roles:` field in frontmatter |
| Inputs not appearing | Check YAML indentation (2 spaces) and that input types are valid, then recompile |
| Wrong branch context | Use `--ref branch-name` in CLI, or select the branch in the GitHub UI dropdown |

## Related Documentation

- [Manual Workflows Example](/gh-aw/examples/manual/) - Example manual workflows
- [Triggers Reference](/gh-aw/reference/triggers/) - Complete trigger syntax including workflow_dispatch
- [TrialOps](/gh-aw/patterns/trial-ops/) - Testing workflows in isolation
- [CLI Commands](/gh-aw/setup/cli/) - Complete gh aw run command reference
- [Templating](/gh-aw/reference/templating/) - Using expressions and conditionals
- [Security Best Practices](/gh-aw/introduction/architecture/) - Securing workflow execution
- [Quick Start](/gh-aw/setup/quick-start/) - Getting started with agentic workflows
