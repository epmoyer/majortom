#/bin/bash

do_build () {
    USE_GOOS=$1
    USE_GOARCH=$2
    TARGET_NAME=$3
    IMAGE_TYPE=$4

    TARGET_DIR=dist/$3
    TARGET_ARCHIVE=dist/images/$3.zip

    echo "Building for $USE_GOOS:$USE_GOARCH into $TARGET_DIR..."

    # Build executable
    GOOS=$USE_GOOS GOARCH=$USE_GOARCH go build -o $TARGET_DIR/majortom

    # Copy supporting files
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/helper.sh $TARGET_DIR

    # Build release
    echo "Building compressed image $TARGET_ARCHIVE..."
    if [ "$IMAGE_TYPE" = "zip" ]; then
        zip -vrj -FS $TARGET_ARCHIVE $TARGET_DIR -x "*.gitkeep" -x "*.DS_Store" > /dev/null
    else
        echo "(TAR NOT YET IMPLEMENTED)"
    fi

}

do_build darwin amd64 macos.amd64 zip
do_build linux amd64 linux.amd64 tar

echo "Done."
