#!/usr/bin/env bash
mkdir -p bin
~/go/bin/gox --osarch="linux/amd64" -output "bin/linux64/{{.Dir}}" ../cmd/...
~/go/bin/gox --osarch="darwin/amd64" -output "bin/darwin64/{{.Dir}}" ../cmd/...
