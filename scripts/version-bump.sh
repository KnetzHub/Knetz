#!/bin/bash
set -e

# Manual version bump script
# Usage: ./scripts/version-bump.sh [major|minor|patch]

BUMP_TYPE=${1:-patch}

echo "🔄 Version Bump Script"
echo "======================"
echo ""

# Get current version
if [ -f VERSION ]; then
  CURRENT_VERSION=$(cat VERSION)
else
  CURRENT_VERSION="0.0.0"
fi

echo "Current version: $CURRENT_VERSION"

# Parse version
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"

# Bump version based on type
case $BUMP_TYPE in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
  *)
    echo "❌ Invalid bump type: $BUMP_TYPE"
    echo "Usage: $0 [major|minor|patch]"
    exit 1
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo "New version: $NEW_VERSION"
echo ""

# Update VERSION file
echo "$NEW_VERSION" > VERSION

# Update README if version mentioned
if [ -f README.md ]; then
  sed -i.bak "s/Version: .*/Version: v$NEW_VERSION/g" README.md && rm README.md.bak
fi

echo "✅ Version updated to $NEW_VERSION"
echo ""
echo "Next steps:"
echo "1. Review changes: git diff"
echo "2. Generate changelog: git-chglog --next-tag v$NEW_VERSION -o CHANGELOG.md"
echo "3. Commit changes: git add . && git commit -m 'chore(release): v$NEW_VERSION'"
echo "4. Create tag: git tag -a v$NEW_VERSION -m 'Release v$NEW_VERSION'"
echo "5. Push: git push origin main && git push origin v$NEW_VERSION"

