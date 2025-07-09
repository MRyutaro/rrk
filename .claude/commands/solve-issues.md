---
name: solve-issues
description: GithubのIssueを解消します．
---

#!/bin/bash
set -e

# Gitリポジトリのルートディレクトリを取得
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$REPO_ROOT" ]; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# fzfがインストールされているか確認
if ! command -v fzf &>/dev/null; then
    echo "❌ Error: 'fzf' is required but not installed. Install it with 'brew install fzf' or similar."
    exit 1
fi

echo "🔍 Fetching open issues..."
echo ""

# オープンなIssueを取得
ISSUES_JSON=$(gh issue list --state open --limit 20 --json number,title,labels)
if [ "$(echo "$ISSUES_JSON" | jq length)" -eq 0 ]; then
    echo "No open issues found."
    exit 0
fi

# コマンドライン引数からIssue番号を取得
if [[ "$1" =~ ^#?[0-9]+$ ]]; then
    ISSUE_NUMBER="${1#\#}"
    # 指定されたIssue番号が存在するか確認
    if ! echo "$ISSUES_JSON" | jq -e ".[] | select(.number == $ISSUE_NUMBER)" > /dev/null; then
        echo "❌ Issue #$ISSUE_NUMBER not found in the list above."
        exit 1
    fi
else
    # fzfでIssueを選択
    SELECTED=$(echo "$ISSUES_JSON" | jq -r '.[] | "\(.number): \(.title)"' | fzf --prompt "🎯 Select an issue to solve: ")
    if [ -z "$SELECTED" ]; then
        echo "❌ No issue selected."
        exit 1
    fi
    ISSUE_NUMBER=$(echo "$SELECTED" | cut -d':' -f1 | tr -d ' ')
fi

# Issueの詳細を取得
echo ""
echo "📖 Fetching issue details..."
ISSUE_DATA=$(gh issue view $ISSUE_NUMBER --json title,body,labels,assignees)
ISSUE_TITLE=$(echo "$ISSUE_DATA" | jq -r '.title')
ISSUE_BODY=$(echo "$ISSUE_DATA" | jq -r '.body // ""')
LABELS=$(echo "$ISSUE_DATA" | jq -r '.labels[].name' | tr '\n' ',' | sed 's/,$//')

echo "🎯 Selected Issue #$ISSUE_NUMBER: $ISSUE_TITLE"
if [ -n "$LABELS" ]; then
    echo "🏷️  Labels: $LABELS"
fi
echo ""

# ブランチ作成
BRANCH_NAME="fix/issue-$ISSUE_NUMBER"
echo "🌿 Creating branch: $BRANCH_NAME"
if git show-ref --verify --quiet refs/heads/"$BRANCH_NAME"; then
    echo "⚠️  Branch $BRANCH_NAME already exists. Switching to it."
    git checkout "$BRANCH_NAME"
else
    git checkout -b "$BRANCH_NAME"
fi
echo ""

# プロジェクト名を取得
PROJECT_NAME=$(basename "$REPO_ROOT")

# Claude AIにIssue解決を依頼
echo "🤖 Preparing context for Claude AI..."
CLAUDE_PROMPT="# 🎯 Issue #$ISSUE_NUMBER: $ISSUE_TITLE

## 📝 Issue Description
$ISSUE_BODY

## 🏷️ Labels
$LABELS

## 📁 Project Context
This is the '$PROJECT_NAME' project.

### Project Structure:
\`\`\`
$(cd "$REPO_ROOT" && find . -type f \( -name "*.go" -o -name "*.py" -o -name "*.js" -o -name "*.ts" -o -name "*.md" -o -name "Makefile" -o -name "go.mod" -o -name "package.json" \) | head -20)
\`\`\`

## 🎯 Task
Please analyze this issue and provide a complete solution that:
1. ✅ Follows best practices and project conventions
2. 🧪 Includes tests if needed
3. 📚 Updates documentation if necessary
4. 🔧 Implements the requested feature/fix

After implementation:
- I'll run tests with \`make test\` or equivalent
- I'll run linting with \`make lint\` or equivalent
- I'll commit and create a PR

Let's solve this issue step by step!"

echo "$CLAUDE_PROMPT"
echo ""
echo "🚀 Now implementing the solution..."
echo "ℹ️  You can start working on the issue. When done:"
echo "   1. Run: git add ."
echo "   2. Run: git commit -m 'fix: resolve issue #$ISSUE_NUMBER'"
echo "   3. Run: git push origin $BRANCH_NAME"
echo "   4. Run: gh pr create --title 'Fix #$ISSUE_NUMBER: $ISSUE_TITLE' --body 'Closes #$ISSUE_NUMBER'"
