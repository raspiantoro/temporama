APP_NAME	= temporama
VERSION		?= $(shell git describe --always --tags)
GIT_COMMIT	?= $(shell git rev-parse --short HEAD)
BUILD_DATE	= $(shell date '+%Y-%m-%d-%H:%M:%S')

install:
	@echo "Install project dependencies"
	go get -v ./...

build: install
	@echo "Building ${APP_NAME} $(VERSION) $(GIT_COMMIT)"
	go build -ldflags "-X github.com/raspiantoro/temporama/info.Server=${APP_NAME} -X github.com/raspiantoro/temporama/info.Version=${VERSION} -X github.com/raspiantoro/temporama/info.GitCommit=${GIT_COMMIT} -X github.com/raspiantoro/temporama/info.BuildDate=${BUILD_DATE}" -o bin/${APP_NAME}

run: build
	@echo "Running ${APP_NAME} $(VERSION) $(GIT_COMMIT)"
	bin/${APP_NAME} ${ARGS}
