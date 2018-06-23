BINARY_NAME=snip
LDFLAGS=-X main.date $(date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.gitCommit $(git rev-parse HEAD)
DOCKER_REPO=saywhat1/snip
VERSION=$(git rev-parse --short=5 HEAD)

all: test build
build:
	go build -o dist/${BINARY_NAME} -v .
test:
	go test -v ./...
clean:
	dep ensure
	rm -f dist/*
run-dev: build
	docker-compose up -d
	./dist/$(BINARY_NAME) --config example-config.yml
stop:
	docker-compose down

docker-build:
	docker build --build-arg LDFLAGS=${LDFLAGS} -t ${DOCKER_REPO}:${VERSION} -t ${DOCKER_REPO}:latest .
docker-push:
	docker push ${DOCKER_REPO}
publish: docker-build docker-push
