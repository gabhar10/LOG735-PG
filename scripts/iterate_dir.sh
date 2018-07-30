#!/bin/bash

set -e 

cd ../ft
for dir in $(pwd)/*; do pushd $dir; go test; popd; done
