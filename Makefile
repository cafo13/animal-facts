CGO_ENABLED=0
GOOS=linux

test:
	go test ./... --tags integration

internal-api-generate-swagger:
	swag init --generalInfo server.go --dir internal-api/server/,internal-api/api/,internal-api/handler/ --output internal-api/docs/

internal-api-run:
	go run cmd/internal-api/main.go

internal-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-internal-api cmd/internal-api/main.go

public-api-generate-swagger:
	swag init --generalInfo server.go --dir public-api/server/,public-api/api/,public-api/handler/ --output public-api/docs/

public-api-run:
	go run cmd/public-api/main.go

public-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-public-api cmd/public-api/main.go
