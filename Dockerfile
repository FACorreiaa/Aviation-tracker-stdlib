FROM node:latest as assets
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
COPY postcss.config.cjs ./
COPY fonts.css ./
RUN mkdir -p controller/static/css controller/static/fonts
RUN npm install --ci
RUN npm run fonts

##CONFIGURE AIR
FROM golang:1.21 as base
LABEL maintainer="a11199"
LABEL description="Base image for building Go applications with Air and Delve."
FROM base as dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /app
CMD ["air"]

### CONFIGURE DEBUG
FROM dev as debug
LABEL maintainer="a11199"
LABEL description="Base image for building Go applications with Air and Delve."
WORKDIR /app
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY . .
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -gcflags "all=-N -l" -o /stay-healthy-backend ./*.go
CMD ["dlv", "--listen=127.0.0.1:40000", "--headless=true", "--api-version=2", "exec", "--accept-multiclient",  "/stay-healthy-backend"]

FROM golang:latest
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
COPY --from=assets /app/controller/static/css/* ./controller/static/css/
COPY --from=assets /app/controller/static/fonts/* ./controller/static/fonts/
RUN CGO_ENABLED=0 go build -o /app/server
EXPOSE 6969
ENTRYPOINT ["/app/server"]
