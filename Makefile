install-mockgen:
	go install github.com/golang/mock/mockgen@latest

FACTS_API_DIR = ./backend/facts-api
generate-mocks: install-mockgen $(FACTS_API_DIR)/*/*.go
	for file in $^ ; do \
		echo "Hello" $${file} ; \
	done

run-backend-tests: generate-mocks
	go test ./backend/facts-api/...
