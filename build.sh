#/bin/bash

do_build () {
    GOOS=$1
    GOARCH=$2
    TARGET_DIR=$3

    echo "Building for $GOOS:$GOARCH into $TARGET_DIR..."

    go build -o $TARGET_DIR/majortom
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/helper.sh $TARGET_DIR
}

do_build darwin amd64 dist/darwin.amd64
do_build linux amd64 dist/linux.amd64

echo "Done."
# TARGET_DIR=dist/darwin.amd64
# GOOS=darwin GOARCH=amd64 go build -o $TARGET_DIR/majortom
# copy_supporting_files $TARGET_DIR

# TARGET_DIR=dist/linux.amd64
# GOOS=linux GOARCH=amd64 go build -o $TARGET_DIR/majortom
# copy_supporting_files $TARGET_DIR
