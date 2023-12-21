genswagger:
	swagger generate server -f ./swagger.yaml -A bundle-service --default-scheme=http

build: build-server build-bundler

build-server:
	go build -o build/bundle-service-server ./cmd/bundle-service-server/main.go

build-bundler:
	go build -o build/bundler ./cmd/bundler/main.go