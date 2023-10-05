#!/bin/bash

# I don't even ..
go env -w GOFLAGS="-buildvcs=false"

# Install tools to test our code-quality.
go install golang.org/x/lint/golint@master
go install honnef.co/go/tools/cmd/staticcheck@master

# Run the static-check tool
for i in */; do
    echo "Running golint on $i"
    cd $i
    t=$(mktemp)
    staticcheck -checks all ./... > $t
    if [ -s $t ]; then
        echo "Found errors via 'staticcheck'"
        cat $t
        rm $t
        exit 1
    fi
    rm $t
    cd ..
done



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
