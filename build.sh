ENV="PROD"
VERSION="1.0.0"

go build -o out/ -ldflags "-X main.Environment=${ENV} -X main.Version=${VERSION} -w -s -H=windowsgui" .
cp ./config.yaml ./out/config.yaml
