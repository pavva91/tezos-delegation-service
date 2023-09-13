#!/bin/bash
# Create symlink
# ln -s ~/go/src/github.com/pavva91/kiln-tezos-delegation-service/githooks/pre-commit.sh ~/go/src/github.com/pavva91/kiln-tezos-delegation-service/.git/hooks/pre-commit

echo "Pre-push Git Hook: Run test suite, abort push if test fail"

# go test ./...
# TODO: cd $(mod_root $FILE) || exit; go test [$ARGS] ./...
