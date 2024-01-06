CGO_ENABLED=0
GOOS=linux

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOTPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SWAG?=$(GOBIN)/swag

.PHONY: $(SWAG)
$(SWAG): $(GOBIN)
	test -s $(GOBIN)/swag || go install github.com/swaggo/swag/cmd/swag@latest

test:
	go test ./... --tags integration

internal-api-generate-swagger: $(SWAG)
	swag init --generalInfo server.go --dir internal-api/server/,internal-api/api/,internal-api/handler/ --output internal-api/docs/

internal-api-run:
	go run cmd/internal-api/main.go

internal-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-internal-api cmd/internal-api/main.go

public-api-generate-swagger: $(SWAG)
	swag init --generalInfo server.go --dir public-api/server/,public-api/api/,public-api/handler/ --output public-api/docs/

public-api-run:
	go run cmd/public-api/main.go

public-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-public-api cmd/public-api/main.go

prepare-release-version:
	sed -i "s/@version         .*\..*\..*/@version         $(VERSION)/g" public-api/server/server.go
	sed -i 's/"starting public animal facts api .*\..*\..*"/"starting public animal facts api $(VERSION)"/g' public-api/server/server.go
	sed -i "s/@version         .*\..*\..*/@version         $(VERSION)/g" internal-api/server/server.go
	sed -i 's/"starting internal animal facts api .*\..*\..*"/"starting internal animal facts api $(VERSION)"/g' internal-api/server/server.go
	echo $(VERSION) > version.txt

release-version: prepare-release-version public-api-generate-swagger internal-api-generate-swagger
