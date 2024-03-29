BIN_DIR=$(CURDIR)/bin

# install dependencies
install-deps:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN_DIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0

# build all binaries
build:
	go build -o bin cmd/migrator/migrator.go

# generate go code from proto files
generate:
	make generate-sso-protobuf

# generate protobuf for sso application
generate-sso-protobuf:
	mkdir -p gen/protobuf/v1/sso
	protoc --proto_path api/proto/v1/sso \
	--go_out=gen/protobuf/v1/sso \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=gen/protobuf/v1/sso \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/proto/v1/sso/sso.proto

# apply all database migrations
migrate-up:
	./bin/migrator --config=config/app.yaml --t=up

# discard all database migrations
migrate-down:
	./bin/migrator --config=config/app.yaml --t=down