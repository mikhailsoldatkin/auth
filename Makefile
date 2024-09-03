include .env
LOCAL_BIN:=$(CURDIR)/bin
USER_V1:=user_v1
AUTH_V1:=auth_v1
ACCESS_V1:=access_v1
REPO:=github.com/mikhailsoldatkin/auth
CERT_FOLDER:=cert

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.21.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.21.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/gojuno/minimock/v3


generate:
	mkdir -p pkg/swagger
	make generate-user-api
	make generate-auth-api
	make generate-access-api
	$(LOCAL_BIN)/statik -f -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-user-api:
	mkdir -p pkg/$(USER_V1)
	protoc --proto_path=api/$(USER_V1) --proto_path=vendor.protogen \
	--go_out=pkg/$(USER_V1) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/$(USER_V1) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out=lang=go:pkg/$(USER_V1) --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_out=pkg/$(USER_V1) --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/$(USER_V1)/user.proto

generate-auth-api:
	mkdir -p pkg/$(AUTH_V1)
	protoc --proto_path api/$(AUTH_V1) \
	--go_out=pkg/$(AUTH_V1) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/$(AUTH_V1) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/$(AUTH_V1)/auth.proto

generate-access-api:
	mkdir -p pkg/$(ACCESS_V1)
	protoc --proto_path api/$(ACCESS_V1) \
	--go_out=pkg/$(ACCESS_V1) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/$(ACCESS_V1) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/$(ACCESS_V1)/access.proto

local-migrations-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} status -v

local-migrations-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} up -v

local-migrations-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} down -v

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=${REPO}/internal/service/...,${REPO}/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=${REPO}/internal/service/...,${REPO}/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

gen-cert:
	mkdir -p $(CERT_FOLDER)
	openssl genrsa -out $(CERT_FOLDER)/ca.key 4096 && \
    openssl req -new -x509 -key $(CERT_FOLDER)/ca.key -sha256 -subj "/C=RU/ST=Moscow/O=Test, Inc." -days 365 -out $(CERT_FOLDER)/ca.cert && \
    openssl genrsa -out $(CERT_FOLDER)/service.key 4096 && \
    openssl req -new -key $(CERT_FOLDER)/service.key -out $(CERT_FOLDER)/service.csr -config $(CERT_FOLDER)/certificate.conf && \
    openssl x509 -req -in $(CERT_FOLDER)/service.csr -CA $(CERT_FOLDER)/ca.cert -CAkey $(CERT_FOLDER)/ca.key -CAcreateserial \
        -out $(CERT_FOLDER)/service.pem -days 365 -sha256 -extfile $(CERT_FOLDER)/certificate.conf -extensions req_ext
