LOCAL_BIN:=$(CURDIR)/bin
CLIENT_APP_NAME:=database_client

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

test:
	go clean -testcache
	go test -v ./... -count=1

build-client:
	go build -o $(LOCAL_BIN)/$(CLIENT_APP_NAME) cmd/client/main.go
