---
name: solve-issues
description: GithubのIssueを解消します．
---

set -e

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$REPO_ROOT" ]; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

echo "🔍 Fetching open issues..."
echo ""

# オープンなIssueを取得して表示
ISSUES=$(gh issue list --state open --limit 20 --json number,title,labels)
if [ "$(echo "$ISSUES" | jq length)" -eq 0 ]; then
    echo "No open issues found."
    exit 0
fi

echo "📋 Open Issues:"
echo "$ISSUES" | jq -r '.[] | "#\(.number) \(.title)"' | nl -v0 -s". "
echo ""

# ユーザーにIssue選択を促す
while true; do
    read -p "🎯 Enter issue number to solve: #" ISSUE_NUMBER
    if [[ "$ISSUE_NUMBER" =~ ^[0-9]+$ ]]; then
        # Issue番号が存在するかチェック
        if echo "$ISSUES" | jq -e ".[] | select(.number == $ISSUE_NUMBER)" > /dev/null; then
            break
        else
            echo "❌ Issue #$ISSUE_NUMBER not found in the list above."
        fi
    else
        echo "❌ Please enter a valid issue number."
    fi
done

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

# Claude AIにIssue解決を依頼
echo "🤖 Preparing context for Claude AI..."
CLAUDE_PROMPT="# 🎯 Issue #$ISSUE_NUMBER: $ISSUE_TITLE

## 📝 Issue Description
$ISSUE_BODY

## 🏷️ Labels
$LABELS

## 📁 Project Context
This is the 'rrk' project - a Go-based CLI tool for enhanced shell history management.

### Project Structure:
\`\`\`
$(cd "$REPO_ROOT" && find . -type f -name "*.go" -o -name "*.md" -o -name "Makefile" -o -name "go.mod" | head -20)
\`\`\`

## 🎯 Task
Please analyze this issue and provide a complete solution that:
1. ✅ Follows Go best practices and project conventions
2. 🧪 Includes tests if needed
3. 📚 Updates documentation if necessary
4. 🔧 Implements the requested feature/fix

After implementation:
- I'll run tests with \`make test\`
- I'll run linting with \`make lint\`
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
