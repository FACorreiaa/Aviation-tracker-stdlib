##CONFIGURE AIR
FROM golang:1.21 as base

LABEL maintainer="a11199"
LABEL description="Base image for building Go applications with Air and Delve."

## Create another stage called "dev" that is based off of our "base" stage (so we have golang available to us)
FROM base as dev

## Install the air binary so we get live code-reloading when we save files
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run the air command in the directory where our code will live
WORKDIR /opt/app/api
CMD ["air"]

### CONFIGURE DEBUG
FROM dev as debug

LABEL maintainer="a11199"
LABEL description="Base image for building Go applications with Air and Delve."

WORKDIR /opt/app/api
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY . .
COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -gcflags "all=-N -l" -o /stay-healthy-backend ./*.go
CMD ["dlv", "--listen=127.0.0.1:40000", "--headless=true", "--api-version=2", "exec", "--accept-multiclient",  "/stay-healthy-backend"]

### MAIN
FROM debug as built

LABEL maintainer="a11199"
LABEL description="Base image for building Go applications with Air and Delve."

WORKDIR /go/app/api

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./... && \
    go build -o /tmp/stay-healthy-backend ./*.go

ENTRYPOINT ["/tmp/stay-healthy-backend"]

FROM busybox

COPY --from=built /tmp/stay-healthy-backend /usr/bin/stay-healthy-backend
CMD ["stay-healthy-backend", "start"]

#"--security-opt='apparmor=unconfined'", "--cap-add=SYS_PTRACE"
