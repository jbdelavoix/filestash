all:
	make build_frontend_all
	make build_backend_all

build_frontend_all:
	make build_frontend_init
	make build_frontend

build_frontend_init:
	docker run -v $(shell pwd):/app -w /app node:13-alpine /bin/ash -c " \
		apk add git && \
		rm -rf node_modules && \
		npm install --silent \
	"

build_frontend:
	docker run -v $(shell pwd):/app -w /app node:13-alpine /bin/ash -c " \
		NODE_ENV=production npm run build \
	"

build_backend_all:
	make build_backend_init
	make build_backend

build_backend_init:
	docker run -v $(shell pwd):/app -w /app golang:1.19-buster /bin/bash -c " \
		apt-get update > /dev/null && \
		apt-get install -y libvips-dev curl > /dev/null 2>&1 && \
		go generate -x ./server/... \
	"

build_backend:
	docker run -v $(shell pwd):/app -w /app golang:1.19-buster /bin/bash -c " \
		apt-get update > /dev/null && \
		apt-get install -y libvips-dev curl > /dev/null 2>&1 && \
		CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -mod=vendor --tags "fts5" -o dist/filestash server/main.go && \
		mkdir -p ./dist/data/state/config/ && \
		cp config/config.json ./dist/data/state/config/config.json \
	"

serve:
	docker run -p 8334:8334 -v $(shell pwd)/dist:/app -w /app debian:stable-slim /app/filestash

go_fmt:
	docker run -v $(shell pwd):/app -w /app golang:1.19-buster /bin/bash -c " \
		go fmt ./server/... \
	"
