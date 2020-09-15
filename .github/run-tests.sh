#!/bin/bash

# Install tools to test our code-quality.
go get -u golang.org/x/lint/golint
go get -u honnef.co/go/tools/cmd/staticcheck

# Run the static-check tool
t=$(mktemp)
staticcheck -checks all ./... > $t
if [ -s $t ]; then
    echo "Found errors via 'staticcheck'"
    cat $t
    rm $t
    exit 1
fi
rm $t



# At this point failures cause aborts
set -e

# Run the linter-tool
for i in */; do
    echo "Running golint on $i"
    cd $i
    golint -set_exit_status ./...
    cd ..
done

# Run the vet-tool
for i in */; do
    echo "Running go vet $i"
    cd $i
    go vet ./...
    cd ..
done


# Run our golang tests
cd foth && go test ./...
