#/bin/bash

do_build () {
    USE_GOOS=$1
    USE_GOARCH=$2
    TARGET_DIR=$3

    echo "Building for $USE_GOOS:$USE_GOARCH into $TARGET_DIR..."

    GOOS=$USE_GOOS GOARCH=$USE_GOARCH go build -o $TARGET_DIR/majortom
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/helper.sh $TARGET_DIR
}

do_build darwin amd64 dist/darwin.amd64
do_build linux amd64 dist/linux.amd64

echo "Done."
