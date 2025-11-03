#!/usr/bin/env python3
"""
Script to create GitHub issues from ISSUES.md
Updates references to match ARCHITECTURE-V3.md (YugabyteDB, DragonflyDB, ScyllaDB, RunPod)

Usage:
    # Dry run (no issues created)
    python3 create_github_issues.py --dry-run
    
    # Create all issues
    python3 create_github_issues.py
    
    # Create specific issue by number (e.g., ISSUE-001 to ISSUE-005)
    python3 create_github_issues.py --range 1-5
"""

import re
import json
import subprocess
import sys
import argparse
from pathlib import Path

# Tech stack replacements based on ARCHITECTURE-V3.md
REPLACEMENTS = {
    "PostgreSQL": "YugabyteDB",
    "postgres": "yugabyte",
    "postgresql": "yugabyte",
    "Redis": "DragonflyDB",
    "redis": "dragonfly",
    "Aerospike": "ScyllaDB",
    "aerospike": "scylla",
    # Keep RabbitMQ for notification service
    # MongoDB -> Keep for analytics/events (but also mention YugabyteDB for structured data)
}

# Labels mapping
LABEL_MAPPING = {
    "infrastructure": "infrastructure",
    "setup": "setup",
    "docker": "docker",
    "database": "database",
    "library": "backend",
    "backend": "backend",
    "service": "service",
    "auth": "auth",
    "user": "user",
    "content": "content",
    "search": "search",
    "streaming": "streaming",
    "transcoding": "transcoding",
    "cdn": "infrastructure",
    "payment": "payment",
    "ads": "ads",
    "analytics": "analytics",
    "ai": "ai",
    "ml": "ml",
    "notifications": "notifications",
    "websocket": "websocket",
    "app": "app",
    "frontend": "frontend",
    "web": "web",
    "mobile": "mobile",
    "react-native": "mobile",
    "tv": "tv",
    "android-tv": "tv",
    "roku": "tv",
    "tizen": "tv",
    "webos": "tv",
    "fire-tv": "tv",
    "ios": "ios",
    "android": "android",
    "deployment": "deployment",
    "kubernetes": "infrastructure",
    "terraform": "infrastructure",
    "devops": "devops",
    "ci-cd": "ci-cd",
    "monitoring": "monitoring",
    "security": "security",
    "testing": "testing",
    "e2e": "testing",
    "performance": "testing",
    "documentation": "documentation",
    "api": "documentation",
    "architecture": "documentation",
    "operations": "documentation",
    "onboarding": "documentation",
    "production": "production",
    "launch": "launch",
    "admin": "admin",
}

PRIORITY_MAPPING = {
    "P0": "Critical",
    "P1": "High",
    "P2": "Medium",
    "P3": "Low",
}


def apply_replacements(text):
    """Apply tech stack replacements"""
    for old, new in REPLACEMENTS.items():
        text = text.replace(old, new)
    return text


def parse_issues_md(file_path):
    """Parse ISSUES.md and extract issue information"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    issues = []
    
    # Split by issue markers
    # Pattern: ### ISSUE-XXX: Title
    issue_pattern = r'### (ISSUE-\d+): (.+?)(?=\n### |$)'
    
    matches = re.finditer(issue_pattern, content, re.DOTALL)
    
    for match in matches:
        issue_num = match.group(1)
        rest = match.group(2)
        
        # Extract fields
        priority_match = re.search(r'\*\*Priority\*\*:\s*(\w+)', rest)
        priority = priority_match.group(1) if priority_match else "P2"
        
        labels_match = re.search(r'\*\*Labels\*\*:\s*`(.+?)`', rest)
        labels_str = labels_match.group(1) if labels_match else ""
        labels = [label.strip() for label in labels_str.split(',') if label.strip()]
        
        estimate_match = re.search(r'\*\*Estimate\*\*:\s*(.+?)(?=\n|$)', rest)
        estimate = estimate_match.group(1).strip() if estimate_match else "Not specified"
        
        deps_match = re.search(r'\*\*Dependencies\*\*:\s*(.+?)(?=\n\*\*Description|$)', rest, re.DOTALL)
        dependencies = deps_match.group(1).strip() if deps_match else "None"
        
        desc_match = re.search(r'\*\*Description\*\*:\s*\n(.+?)(?=\n\*\*Tasks:|$)', rest, re.DOTALL)
        description = desc_match.group(1).strip() if desc_match else ""
        
        tasks_match = re.search(r'\*\*Tasks\*\*:\s*\n(.+?)(?=\n\*\*Deliverables:|$)', rest, re.DOTALL)
        tasks = tasks_match.group(1).strip() if tasks_match else ""
        
        deliverables_match = re.search(r'\*\*Deliverables\*\*:\s*\n(.+?)(?=\n---|\n##|$)', rest, re.DOTALL)
        deliverables = deliverables_match.group(1).strip() if deliverables_match else ""
        
        # Extract title from first line after issue number
        title = rest.split('\n')[0].strip()
        
        # Apply replacements
        description = apply_replacements(description)
        tasks = apply_replacements(tasks)
        deliverables = apply_replacements(deliverables)
        
        # Build issue body
        body = f"""## Description

{description}

## Priority

**{priority}** ({PRIORITY_MAPPING.get(priority, priority)})

## Estimate

{estimate}

## Dependencies

{dependencies}

## Tasks

{tasks}

## Deliverables

{deliverables}

---
*This issue was auto-generated from ISSUES.md and updated to use YugabyteDB, DragonflyDB, ScyllaDB, and RunPod as per ARCHITECTURE-V3.md*
"""
        
        github_labels = [LABEL_MAPPING.get(label, label) for label in labels if label]
        
        issues.append({
            "number": issue_num,
            "title": title,
            "body": body,
            "labels": list(set(github_labels)),  # Remove duplicates
            "priority": priority,
            "estimate": estimate,
            "dependencies": dependencies,
        })
    
    return sorted(issues, key=lambda x: int(x["number"].split("-")[1]))


def ensure_labels_exist(repo, labels):
    """Ensure all labels exist in the repository"""
    # Get existing labels
    result = subprocess.run(
        ["gh", "label", "list", "--repo", repo, "--json", "name"],
        capture_output=True,
        text=True
    )
    
    existing_labels = set()
    if result.returncode == 0:
        existing_labels_json = json.loads(result.stdout)
        existing_labels = {label["name"] for label in existing_labels_json}
    
    # Create missing labels with colors
    label_colors = {
        "infrastructure": "0052CC",
        "backend": "E6F7FF",
        "service": "7057FF",
        "app": "FBCA04",
        "frontend": "FBCA04",
        "mobile": "FBCA04",
        "web": "FBCA04",
        "tv": "FBCA04",
        "auth": "D4C5F9",
        "user": "D4C5F9",
        "content": "D4C5F9",
        "streaming": "D4C5F9",
        "transcoding": "D4C5F9",
        "payment": "D4C5F9",
        "analytics": "D4C5F9",
        "testing": "EDEDED",
        "documentation": "BFD4F2",
        "devops": "E6F7FF",
        "ci-cd": "E6F7FF",
        "security": "B60205",
        "production": "D73A4A",
        "launch": "D73A4A",
    }
    
    all_labels = set()
    for issue_labels in labels:
        all_labels.update(issue_labels)
    
    for label in all_labels:
        if label not in existing_labels:
            color = label_colors.get(label, "0E8A16")
            result = subprocess.run(
                ["gh", "label", "create", label, "--color", color, "--repo", repo],
                capture_output=True
            )
            if result.returncode == 0:
                print(f"  ✓ Created label: {label}")
            else:
                print(f"  ✗ Failed to create label {label}: {result.stderr}")


def create_issue(repo, issue_data, dry_run=False):
    """Create a GitHub issue"""
    title = issue_data["title"]
    body = issue_data["body"]
    labels = issue_data["labels"]
    
    if dry_run:
        print(f"\n[DRY RUN] Would create issue:")
        print(f"  Title: {title}")
        print(f"  Labels: {', '.join(labels)}")
        print(f"  Priority: {issue_data['priority']}")
        return None
    
    # Create issue with gh CLI
    label_args = []
    for label in labels:
        label_args.extend(["--label", label])
    
    # Add priority as label
    priority_label = f"priority-{issue_data['priority'].lower()}"
    label_args.extend(["--label", priority_label])
    
    # Build command with --body flag instead of stdin
    cmd = ["gh", "issue", "create", "--repo", repo, "--title", title, "--body", body] + label_args
    
    result = subprocess.run(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    if result.returncode == 0:
        issue_url = result.stdout.strip()
        return issue_url
    else:
        print(f"  ✗ Error: {result.stderr.strip()}")
        return None


def main():
    parser = argparse.ArgumentParser(description='Create GitHub issues from ISSUES.md')
    parser.add_argument('--dry-run', '-d', action='store_true', 
                       help='Dry run mode - show what would be created without creating issues')
    parser.add_argument('--range', '-r', type=str,
                       help='Issue range to create (e.g., "1-5" or "1,3,5")')
    parser.add_argument('--repo', type=str,
                       default='abiolaogu/https-github.com-abiolaogu-Video-Streaming_AI-Studio',
                       help='GitHub repository (owner/repo)')
    args = parser.parse_args()
    
    repo = args.repo
    
    # Parse issues
    script_dir = Path(__file__).parent
    issues_file = script_dir / "ISSUES.md"
    
    if not issues_file.exists():
        print(f"ERROR: {issues_file} not found!")
        print("Please ensure ISSUES.md is in the same directory as this script.")
        sys.exit(1)
    
    print(f"Parsing {issues_file}...")
    all_issues = parse_issues_md(issues_file)
    print(f"Found {len(all_issues)} issues")
    
    # Filter issues if range specified
    if args.range:
        if '-' in args.range:
            start, end = map(int, args.range.split('-'))
            issues = [issue for issue in all_issues 
                     if start <= int(issue["number"].split("-")[1]) <= end]
        elif ',' in args.range:
            numbers = [int(n.strip()) for n in args.range.split(',')]
            issues = [issue for issue in all_issues 
                     if int(issue["number"].split("-")[1]) in numbers]
        else:
            issues = all_issues
        print(f"Filtered to {len(issues)} issues")
    else:
        issues = all_issues
    
    if not issues:
        print("No issues to process.")
        sys.exit(0)
    
    if not args.dry_run:
        # Ensure labels exist
        print("\nEnsuring labels exist...")
        all_labels = [issue["labels"] for issue in issues]
        ensure_labels_exist(repo, all_labels)
        
        # Ensure priority labels exist
        priority_labels = [f"priority-{issue['priority'].lower()}" for issue in issues]
        ensure_labels_exist(repo, [priority_labels])
    
    # Create issues
    print(f"\n{'[DRY RUN] ' if args.dry_run else ''}Creating issues...")
    created = []
    failed = []
    
    for i, issue_data in enumerate(issues, 1):
        print(f"\n[{i}/{len(issues)}] {issue_data['number']}: {issue_data['title']}")
        issue_url = create_issue(repo, issue_data, dry_run=args.dry_run)
        
        if issue_url:
            created.append(issue_data['number'])
            if not args.dry_run:
                print(f"  ✓ Created: {issue_url}")
        else:
            failed.append(issue_data['number'])
    
    # Summary
    print("\n" + "="*60)
    print("SUMMARY")
    print("="*60)
    print(f"Total issues processed: {len(issues)}")
    print(f"Created: {len(created)}")
    if failed:
        print(f"Failed: {len(failed)}")
        print(f"Failed issues: {', '.join(failed)}")
    print("\n✓ Done!")
    
    if args.dry_run:
        print("\nTo create the issues, run without --dry-run flag")


if __name__ == "__main__":
    main()
