#/bin/bash
GOOS=darwin GOARCH=amd64 go build -o dist/darwin.amd64/majortom
cp dist/resources/install.sh dist/darwin.amd64/
GOOS=linux GOARCH=amd64 go build -o dist/linux.amd64/majortom
cp dist/resources/install.sh dist/linux.amd64/
