#/bin/bash
GOOS=darwin GOARCH=amd64 go build -o dist/darwin.amd64/majortom
GOOS=linux GOARCH=amd64 go build -o dist/linux.amd64/majortom
