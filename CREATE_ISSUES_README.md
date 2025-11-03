# Creating GitHub Issues from ISSUES.md

This script automatically creates GitHub issues from the `ISSUES.md` file, updating tech stack references to match `ARCHITECTURE-V3.md` (YugabyteDB, DragonflyDB, ScyllaDB, RunPod).

## Prerequisites

1. **GitHub CLI Authentication**
   ```bash
   gh auth login
   ```
   Follow the prompts to authenticate with GitHub.

2. **Verify Authentication**
   ```bash
   gh auth status
   ```

## Usage

### 1. Dry Run (Recommended First)

Test the script without creating any issues:

```bash
python3 create_github_issues.py --dry-run
```

This will:
- Parse all issues from `ISSUES.md`
- Show what would be created
- Not create any actual GitHub issues

### 2. Create All Issues

Once you've verified the dry run looks correct:

```bash
python3 create_github_issues.py
```

This will create all 51 issues in your GitHub repository.

### 3. Create Specific Issues

Create a range of issues:

```bash
python3 create_github_issues.py --range 1-5
```

Create specific issues by number:

```bash
python3 create_github_issues.py --range 1,3,5
```

### 4. Custom Repository

If you need to specify a different repository:

```bash
python3 create_github_issues.py --repo owner/repo-name
```

## What the Script Does

1. **Parses ISSUES.md** - Extracts all 51 issues with their details
2. **Applies Tech Stack Updates**:
   - PostgreSQL → YugabyteDB
   - Redis → DragonflyDB
   - Aerospike → ScyllaDB
   - Updates references to match ARCHITECTURE-V3.md
3. **Creates GitHub Labels** - Automatically creates all necessary labels with appropriate colors
4. **Creates Issues** - Uses GitHub CLI to create issues with:
   - Title from issue name
   - Full description, tasks, deliverables
   - Appropriate labels
   - Priority labels (priority-p0, priority-p1, etc.)

## Tech Stack Updates Applied

Based on ARCHITECTURE-V3.md:

- **PostgreSQL** → **YugabyteDB** (distributed SQL)
- **Redis** → **DragonflyDB** (in-memory cache)
- **Aerospike** → **ScyllaDB** (time-series database)
- **AWS GPU** → **RunPod** (for transcoding - mentioned in transcoding issues)
- **RabbitMQ** → Kept for notification service, but Kafka preferred for main messaging

## Issue Structure

Each issue includes:
- **Description** - What needs to be done
- **Priority** - P0 (Critical), P1 (High), P2 (Medium), P3 (Low)
- **Estimate** - Time estimate in hours
- **Dependencies** - Prerequisite issues
- **Tasks** - Detailed checklist
- **Deliverables** - Expected outputs

## Labels Created

The script automatically creates and uses these labels:

- **Category**: infrastructure, backend, service, app, frontend, mobile, tv
- **Technology**: auth, user, content, streaming, transcoding, payment, analytics
- **Type**: testing, documentation, devops, security, production, launch
- **Priority**: priority-p0, priority-p1, priority-p2, priority-p3

## Troubleshooting

### Authentication Issues

```bash
# Re-authenticate
gh auth login

# Check status
gh auth status
```

### Missing ISSUES.md

Ensure `ISSUES.md` is in the same directory as the script.

### Permission Issues

Make sure you have write access to the repository:
```bash
gh repo view abiolaogu/https-github.com-abiolaogu-Video-Streaming_AI-Studio
```

### Rate Limiting

GitHub has rate limits. If you hit them:
- Wait a few minutes and retry
- Or create issues in batches using `--range`

## Next Steps

After creating issues:

1. **Review Issues** - Check a few issues to ensure they look correct
2. **Set Up Milestones** - Create milestones for each phase:
   - Phase 1: Foundation & Core Infrastructure
   - Phase 2: Authentication & User Management
   - Phase 3: Content Management & Catalog
   - etc.
3. **Create GitHub Project** - Set up a project board to track progress
4. **Link Dependencies** - Use GitHub's dependency feature to link related issues

## Example Output

```
Parsing ISSUES.md...
Found 51 issues

Ensuring labels exist...
  ✓ Created label: infrastructure
  ✓ Created label: backend
  ...

Creating issues...

[1/51] ISSUE-001: Project Setup & Monorepo Structure
  ✓ Created: https://github.com/.../issues/1

[2/51] ISSUE-002: Docker Development Environment
  ✓ Created: https://github.com/.../issues/2

...

SUMMARY
============================================================
Total issues processed: 51
Created: 51
Failed: 0

✓ Done!
```

## Notes

- The script updates database references automatically
- All issues include a note that they were auto-generated and updated per ARCHITECTURE-V3.md
- Issues are created in order (ISSUE-001, ISSUE-002, etc.)
- You can manually edit issues after creation if needed

