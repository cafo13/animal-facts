install-mockgen:
	go install github.com/golang/mock/mockgen@latest

API_DIR = ./backend
generate-mocks: install-mockgen $(API_DIR)/*/*.go
	for file in $^ ; do \
		echo "Hello" $${file} ; \
	done

run-backend-tests: generate-mocks
	go test ./backend/...
