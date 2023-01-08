#/bin/bash
GREEN=$'\033[32m'
CYAN=$'\033[36m'
ENDCOLOR=$'\033[0m'

do_build () {
    USE_GOOS=$1
    USE_GOARCH=$2
    TARGET_NAME=$3
    IMAGE_TYPE=$4

    TARGET_DIR=dist/$3

    echo "${CYAN}$USE_GOOS:$USE_GOARCH${ENDCOLOR} -------------------------------------------"

    echo "Building into $TARGET_DIR..."

    # Build executable
    GOOS=$USE_GOOS GOARCH=$USE_GOARCH go build -o $TARGET_DIR/majortom

    # Copy supporting files
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/helper.sh $TARGET_DIR

    # Build release
    if [ "$IMAGE_TYPE" = "zip" ]; then
        TARGET_ARCHIVE=dist/images/$3.zip
        echo "Building compressed image $TARGET_ARCHIVE..."
        zip -vrj -FS $TARGET_ARCHIVE $TARGET_DIR -x "*.gitkeep" -x "*.DS_Store" > /dev/null
    else
        TARGET_ARCHIVE=dist/images/$3.tgz
        echo "Building compressed image $TARGET_ARCHIVE..."
        #tar -cvzf dist/images/linux.amd64.tgz -C dist/linux.amd64 --exclude=.gitkeep --exclude=.DS_Store .
        tar -cvzf $TARGET_ARCHIVE -C $TARGET_DIR --exclude=.gitkeep --exclude=.DS_Store .
        echo "(TAR NOT YET IMPLEMENTED)"
    fi
    echo ""
}

do_build darwin amd64 macos.amd64 zip
do_build linux amd64 linux.amd64 tar

echo "${GREEN}Done.${ENDCOLOR}"
