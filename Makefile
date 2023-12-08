genswagger:
	swagger generate server -f ./swagger.yaml -A bundle-service --default-scheme=http

build-server:
	go build -o build/bundle-service-server ./cmd/bundle-service-server/main.go

