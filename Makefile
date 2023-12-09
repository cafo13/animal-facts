GOOS=linux
GOARCH=amd64

tests:
	go test ./...

internal-api-generate-swagger:
	swag init --dir cmd/internal-api/ --output internal-api/docs/
#	swag init --dir cmd/internal-api/,internal-api/api/,internal-api/handler/ --output internal-api/docs/

internal-api-run:
	go run cmd/internal-api/main.go

internal-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-internal-api cmd/internal-api/main.go

public-api-generate-swagger:
	swag init --dir cmd/public-api/,public-api/api/,public-api/handler/ --output public-api/docs/

public-api-run:
	go run cmd/public-api/main.go

public-api-build:
	go build -ldflags "-s -w" -o bin/animal-facts-public-api cmd/public-api/main.go
