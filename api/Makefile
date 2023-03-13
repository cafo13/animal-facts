install-mockgen:
	go install github.com/golang/mock/mockgen@latest


API_DIR = .
generate-mocks: install-mockgen $(API_DIR)/*/*
	for file in $^ ; do \
		echo "Hello" $${file} ; \
	done

run-tests: generate-mocks
	go test ./...

install-swaggo:
	go install github.com/swaggo/swag/cmd/swag@latest

generate-docs: install-swaggo
	swag init --parseDependency; rm ./docs/docs.go
