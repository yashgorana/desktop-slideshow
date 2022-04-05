ENV="PROD"
VERSION="1.0.0"
BUILD_DIR="./.out"

rm -rf $BUILD_DIR
go build -o $BUILD_DIR/ -ldflags "-X main.Environment=${ENV} -X main.Version=${VERSION} -w -s -H=windowsgui" .
cp ./config.yaml $BUILD_DIR/config.yaml
