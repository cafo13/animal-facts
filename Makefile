install-mockgen:
	go install github.com/golang/mock/mockgen@latest

API_DIR = ./backend
generate-mocks: install-mockgen $(API_DIR)/*/*.go
	for file in $^ ; do \
		echo "Hello" $${file} ; \
	done

run-backend-tests: generate-mocks
	go test ./backend/...

install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

openapi_http: install-oapi-codegen
	oapi-codegen -generate types -o "backend/facts/openapi/types.gen.go" -package "openapi" "api/openapi/facts.yml"
	oapi-codegen -generate chi-server -o "backend/facts/openapi/api.gen.go" -package "openapi" "api/openapi/facts.yml"
	oapi-codegen -generate types -o "backend/common/client/facts/types.gen.go" -package "facts" "api/openapi/facts.yml"
	oapi-codegen -generate client -o "backend/common/client/facts/client.gen.go" -package "facts" "api/openapi/facts.yml"
