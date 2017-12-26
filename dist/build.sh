#!/bin/bash

declare -a PLATFORMS=(
    'linux/amd64'
    'linux/arm'
)

NAME='tempreader'
VER="${TRAVIS_TAG:1}"
TMP="/tmp/$NAME"
rm -fr "$TMP"

for P in "${PLATFORMS[@]}"; do
    echo "Building $P"

    PTMP="$TMP/$P/$NAME-$VER"
    OS="${P%%/*}"
    ARCH="${P#*/}"
    PKG="$NAME-$VER-$OS-$ARCH.tar.gz"

    mkdir -p "$PTMP"
    GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 go build -ldflags='-s -w' -o "$PTMP/$NAME"

    (
        cp -r "$TRAVIS_BUILD_DIR"/{README.md,LICENSE,contrib/tempreader.{service,default}} "$PTMP"
        tar -C "$PTMP/.." -czf "$TMP/$PKG" ./
        cd "$TMP"
    )
done

cd "$TMP"
sha256sum *.tar.gz | sort -k2 > "$NAME-$VER-checksums-sha256.txt"
