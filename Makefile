BINARY_NAME=cloudin

gitver: BuildVersion=$(git describe --abbrev=0 --tags)
time: BuildTime=${date +'%m-%d-%Y_%H-%M'}

LDFLAGS=--ldflags="-X 'main.BuildTime=${BuildTime} -X 'main.BuildVersion=${BuildVersion}"

build: go build ${LDFLAGS} -o artifact/${BINARY_NAME}-build main.go

compile:
  GOARCH=368 GOOS=linux  go build ${LDFLAGS}  -o artifact/${BINARY_NAME}-linux main.go
  GOARCH=368 GOOS=window go build ${LDFLAGS} -o artifact/${BINARY_NAME}-windows main.go

all: build compile


