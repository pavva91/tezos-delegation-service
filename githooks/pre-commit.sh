#!/bin/bash
# Create symlink
# ln -s ~/go/src/github.com/pavva91/kiln-tezos-delegation-service/githooks/pre-commit.sh ~/go/src/github.com/pavva91/kiln-tezos-delegation-service/.git/hooks/pre-commit

echo "Pre-commit Git Hook: Format staged go files"

STAGED_GO_FILES=$(git diff --cached --name-only -- '*.go')

if $STAGED_GO_FILES '==' ""; then
	echo "No Go Files to Update"
else
	for file in $STAGED_GO_FILES; do
		go fmt "$file"
		git add "$file"
        echo "$file"
	done
fi
