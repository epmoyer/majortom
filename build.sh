#/bin/bash
GREEN=$'\033[32m'
CYAN=$'\033[36m'
ENDCOLOR=$'\033[0m'

# Extract the app name and version numberfrom the main go file by brute force
APP_VERSION=`awk '/^const APP_VERSION/{print $4}' majortom.go | sed -e 's/"//g'`
APP_NAME=`awk '/^const APP_NAME/{print $4}' majortom.go | sed -e 's/"//g'`

echo "Building $APP_NAME version $APP_VERSION..."

do_build () {
    USE_GOOS=$1
    USE_GOARCH=$2
    TARGET_NAME=$3
    IMAGE_TYPE=$4

    TARGET_DIR=dist/builds/$3

    echo "${CYAN}$USE_GOOS:$USE_GOARCH${ENDCOLOR} -------------------------------------------"

    echo "Building into $TARGET_DIR..."

    # Build executable
    GOOS=$USE_GOOS GOARCH=$USE_GOARCH go build -o $TARGET_DIR/majortom

    # Copy supporting files
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/shell_init_snippet.sh $TARGET_DIR

    # Build release
    if [ "$IMAGE_TYPE" = "zip" ]; then
        TARGET_ARCHIVE=dist/images/${APP_NAME}_$APP_VERSION.$TARGET_NAME.zip
        echo "Building compressed image $TARGET_ARCHIVE..."
        zip -vrj -FS $TARGET_ARCHIVE $TARGET_DIR -x "*.gitkeep" -x "*.DS_Store" > /dev/null
    else
        TARGET_ARCHIVE=dist/images/${APP_NAME}_$APP_VERSION.$TARGET_NAME.tgz
        echo "Building compressed image $TARGET_ARCHIVE..."
        tar -czf $TARGET_ARCHIVE -C $TARGET_DIR --exclude=.gitkeep --exclude=.DS_Store .
    fi
    echo ""
}

do_build darwin amd64 macos.amd64 zip
do_build linux amd64 linux.amd64 tar

echo "${GREEN}Done.${ENDCOLOR}"
