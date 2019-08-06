#!/bin/bash

echo "Hello to stdout"
echo "Hello to stderr" >&2
let i=1;
for arg in "$@"; do
    echo "ARG${i} $arg"
    i=$(( i + 1 ))
done

# Return first argument as exit code
exit $1