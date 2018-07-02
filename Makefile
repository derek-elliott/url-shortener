BINARY_NAME=snip
DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
GIT_COMMIT=$(shell git rev-parse HEAD)
LDFLAGS="-X main.date=$(DATE) -X main.gitCommit=$(GIT_COMMIT)"
DOCKER_REPO=saywhat1/snip
VERSION=$(shell git rev-parse --short=5 HEAD)

help: ## This help message
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
all: test build  ## Runs tests and build the binary
build: ## Builds the binary in the dist directory
	go build -o dist/${BINARY_NAME} -v .
test:  ## Run the tests
	go test -v ./...
clean: ## Ensure all dependencies are up to date, clean vendor and remove binaries from dist directory
	dep ensure
	rm -f dist/*
run-dev: build ## Starts dev Postgres and Redis docker Containers and starts the service
	docker-compose up -d
	sleep 5
	./dist/$(BINARY_NAME) --config example-config.yml
stop: ## Stops the dev Postgres and Redis Docker containers
	docker-compose down

docker-build: ## Builds the docker container for the service, tags it with the version and latest
	docker build --build-arg LDFLAGS=${LDFLAGS} -t ${DOCKER_REPO}:$(VERSION) -t ${DOCKER_REPO}:latest .
docker-push: ## Pushes built containers to the Docker repo
	docker push ${DOCKER_REPO}
publish: docker-build docker-push ## Builds the container and pushes to the remote repo
