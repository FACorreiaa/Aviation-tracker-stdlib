project_name = stay-healthy-backend
image_name = stay-healthy-backend

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

run-dev:
	go run main.go

requirements:
	make clean-packages
	go mod tidy

clean-packages:
	go clean -modcache

up:
	make up-silent
	make shell

build:
	docker build -t $(image_name) .

push:
	docker build -t $(image_name) .

build-no-cache:
	docker build --no-cache -t $(image_name) . --security-opt=seccomp:unconfined

up-silent:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./main

up-silent-prefork:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app -prod

delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

shell:
	docker exec -it $(project_name) /bin/sh

compose-up:
	make delete-container-if-exist
	docker compose up -d

stop:
	docker stop $(project_name)

start:
	docker start $(project_name)

swag-init:
	swag init --parseDependency

go-test:
	go test ./...

go-bench:
	go test -bench .

test_coverage:
	go test ./... -coverprofile=coverage.out

run-app:
	docker compose run --rm app air init

run-tidy:
	docker compose run --rm app go mod tidy

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

start-redis:
	redis-server --port 6969
