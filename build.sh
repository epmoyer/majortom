#/bin/bash
GREEN=$'\033[32m'
CYAN=$'\033[36m'
ENDCOLOR=$'\033[0m'

# Extract the app name and version number from the main go file by brute force
APP_VERSION=`awk '/^const APP_VERSION/{print $4}' majortom.go | sed -e 's/"//g'`
APP_NAME=`awk '/^const APP_NAME/{print $4}' majortom.go | sed -e 's/"//g'`

echo "Building $APP_NAME version $APP_VERSION..."

if ! command -v gtar &> /dev/null
then
    TAR_APP=tar
    echo "gtar not found. Will use tar instead."
    echo "NOTE: If you are building on MacOS then the native MacOS tar will cause"
    echo "      (benign) warnings about the extended header keyword LIBARCHIVE.xattr.com.dropbox.attrs"
    echo "      when un-tar-ing the resulting tar archive on Linux systems.  Using gtar avoids this."
else
    TAR_APP=gtar
    echo "gtar found. Will use gtar."
fi
# TAR_APP=tar

do_build () {
    USE_GOOS=$1
    USE_GOARCH=$2
    IMAGE_TYPE=$3

    OS_NAME=$USE_GOOS
    if [[ $OS_NAME = "darwin" ]]
    then
        # We will use the name 'macos' instead of 'darwin' for build images to provide clarity
        # for a broader range of macos users.
        OS_NAME='macos'
    fi

    TARGET_NAME=$OS_NAME.$USE_GOARCH

    BUILD_DIR_BASE=dist/builds
    BUILD_DIR_FINAL=$APP_NAME\_$APP_VERSION.$TARGET_NAME
    TARGET_DIR=$BUILD_DIR_BASE/$BUILD_DIR_FINAL

    echo "${CYAN}$USE_GOOS:$USE_GOARCH${ENDCOLOR} -------------------------------------------"

    echo "Building into: $TARGET_DIR..."

    echo "Making directory (if missing): $TARGET_DIR..."
    mkdir -p $TARGET_DIR

    # Build executable
    GOOS=$USE_GOOS GOARCH=$USE_GOARCH go build -o $TARGET_DIR/majortom

    # Copy supporting files
    cp dist/resources/install.sh $TARGET_DIR
    cp dist/resources/shell_init_snippet.sh $TARGET_DIR

    # I don't really like the non-determinism of this sleep command, but without it the tar below
    # would sometimes throw the warning "gtar: ./majortom: file changed as we read it".
    # Presumably the go build is backgrounding some stage of operation, or the file system is
    # synching.  I tried 'sync' and 'wait', but only `sleep 1` appears to (for now) have reliably
    # squashed the warning (which is really an error to us, if the build has not completed before
    # we tar).
    sleep 1

    # Build release
    if [ "$IMAGE_TYPE" = "zip" ]; then
        TARGET_ARCHIVE=dist/images/${APP_NAME}_$APP_VERSION.$TARGET_NAME.zip
        echo "Building compressed image $TARGET_ARCHIVE..."
        zip -vrj -FS $TARGET_ARCHIVE $TARGET_DIR -x "*.gitkeep" -x "*.DS_Store" > /dev/null
    else
        TARGET_ARCHIVE=dist/images/${APP_NAME}_$APP_VERSION.$TARGET_NAME.tgz
        echo "Building compressed image $TARGET_ARCHIVE..."
        $TAR_APP -czf $TARGET_ARCHIVE -C $BUILD_DIR_BASE --exclude=.gitkeep --exclude=.DS_Store $BUILD_DIR_FINAL
    fi
    echo ""
}

# Run builds for all target platforms
do_build darwin amd64 zip
do_build linux amd64 tar

echo "${GREEN}Done.${ENDCOLOR}"
